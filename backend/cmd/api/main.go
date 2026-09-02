package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PurePras/shri-ram-service/backend/internal/config"
	"github.com/PurePras/shri-ram-service/backend/internal/database"
	"github.com/PurePras/shri-ram-service/backend/internal/handler"
	"github.com/PurePras/shri-ram-service/backend/internal/middleware"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler(pool))
	registerAPIRoutes(mux, pool, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: withCORS(mux, cfg.CORSAllowedOrigin),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	log.Printf("Shri Ram API running on http://localhost:%s", cfg.ServerPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
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

	requireAdmin := middleware.RequireRole(cfg.JWTSecret, model.RoleAdmin)

	categoryHandler.RegisterRoutes(mux, requireAdmin)
	productHandler.RegisterRoutes(mux, requireAdmin)
	deliveryAreaHandler.RegisterRoutes(mux, requireAdmin)
	deliverySlotHandler.RegisterRoutes(mux, requireAdmin)
	userHandler.RegisterRoutes(mux)
	orderHandler.RegisterRoutes(mux, requireAdmin)
	authHandler.RegisterRoutes(mux)
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
