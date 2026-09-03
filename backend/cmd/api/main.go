package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PurePras/shri-ram-service/backend/internal/config"
	"github.com/PurePras/shri-ram-service/backend/internal/database"
	"github.com/PurePras/shri-ram-service/backend/internal/handler"
	"github.com/PurePras/shri-ram-service/backend/internal/middleware"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

// runHealthcheck is invoked as `api -healthcheck` — a Docker HEALTHCHECK exec's a new
// process inside the *same* container, and the distroless runtime image has no shell or
// curl for the usual `CMD curl -f ...`. Re-running this binary against the server it's
// already running alongside (same localhost, same PORT) is the standard workaround.
func runHealthcheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	resp, err := http.Get("http://localhost:" + port + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	resp.Body.Close()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		runHealthcheck()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		// Config (including LogFormat) isn't available yet, so this one line can't use slog.
		panic("load config: " + err.Error())
	}

	logger := newLogger(cfg.LogFormat)
	slog.SetDefault(logger)

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Error("run migrations", "error", err)
		os.Exit(1)
	}

	// `api -migrate` applies migrations and exits — used by CI to prepare a database for
	// the integration suite, and available as the same step in a real deploy pipeline
	// that wants migrating separated from starting the server.
	if len(os.Args) > 1 && os.Args[1] == "-migrate" {
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler(pool))
	// Not reverse-proxied by nginx (see nginx/conf.d/app.conf) — only reachable from
	// other containers on the compose network, e.g. Prometheus. No auth of its own.
	mux.Handle("/metrics", promhttp.Handler())
	registerAPIRoutes(mux, pool, cfg)

	// Built up innermost-first: the general rate limiter sits closest to the mux so a
	// browser's automatic OPTIONS preflight (handled by withCORS) never consumes a token.
	// Metrics wraps the mux directly so r.Pattern (set by the mux's own routing) is
	// already populated by the time it reads it back out.
	generalLimiter := middleware.NewIPRateLimiter(5, 20)
	var root http.Handler = middleware.Metrics(mux)
	root = generalLimiter.Middleware(root)
	root = withCORS(root, cfg.CORSAllowedOrigin)
	root = middleware.SecurityHeaders(root)
	root = middleware.RequestLogger(logger)(root)
	root = middleware.Recover(logger)(root)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	logger.Info("Shri Ram API running", "port", cfg.ServerPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func newLogger(format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func registerAPIRoutes(mux *http.ServeMux, pool *pgxpool.Pool, cfg *config.Config) {
	userRepo := repository.NewUserRepository(pool)

	categoryHandler := handler.NewCategoryHandler(service.NewCategoryService(repository.NewCategoryRepository(pool)))
	productHandler := handler.NewProductHandler(service.NewProductService(repository.NewProductRepository(pool)))
	deliveryAreaHandler := handler.NewDeliveryAreaHandler(service.NewDeliveryAreaService(repository.NewDeliveryAreaRepository(pool)))
	deliverySlotHandler := handler.NewDeliverySlotHandler(service.NewDeliverySlotService(repository.NewDeliverySlotRepository(pool)))
	userHandler := handler.NewUserHandler(service.NewUserService(userRepo))
	orderHandler := handler.NewOrderHandler(service.NewOrderService(repository.NewOrderRepository(pool)))
	authHandler := handler.NewAuthHandler(service.NewAuthService(userRepo, cfg.JWTSecret), cfg.JWTSecret, cfg.CookieSecure)

	requireAuth := middleware.RequireAuth(cfg.JWTSecret)
	requireAdmin := middleware.RequireRole(cfg.JWTSecret, model.RoleAdmin)
	// Stricter than the general per-IP limit — login/register/password-change are the
	// actual brute-force/credential-stuffing targets, not the API as a whole.
	authLimiter := middleware.NewIPRateLimiter(1, 5)

	categoryHandler.RegisterRoutes(mux, requireAdmin)
	productHandler.RegisterRoutes(mux, requireAdmin)
	deliveryAreaHandler.RegisterRoutes(mux, requireAdmin)
	deliverySlotHandler.RegisterRoutes(mux, requireAdmin)
	userHandler.RegisterRoutes(mux)
	orderHandler.RegisterRoutes(mux, requireAuth, requireAdmin)
	authHandler.RegisterRoutes(mux, authLimiter.Guard)
}

func withCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := "ok"
		dbStatus := "ok"
		if err := pool.Ping(r.Context()); err != nil {
			status = "degraded"
			dbStatus = "unreachable"
		}

		response := map[string]string{
			"status":   status,
			"service":  "shri-ram-api",
			"database": dbStatus,
		}

		json.NewEncoder(w).Encode(response)
	}
}
