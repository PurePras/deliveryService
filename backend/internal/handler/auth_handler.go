package handler

import (
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/auth"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/middleware"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

type AuthHandler struct {
	svc          *service.AuthService
	jwtSecret    string
	cookieSecure bool
}

func NewAuthHandler(svc *service.AuthService, jwtSecret string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{svc: svc, jwtSecret: jwtSecret, cookieSecure: cookieSecure}
}

// authLimiter is a stricter, auth-specific rate limit (see middleware.IPRateLimiter.Guard) —
// applied outermost on the sensitive routes, before the cost of checking auth/hitting the DB.
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, authLimiter RouteGuard) {
	mux.HandleFunc("POST /auth/register", authLimiter(h.Register))
	mux.HandleFunc("POST /auth/login", authLimiter(h.Login))
	mux.HandleFunc("POST /auth/logout", h.Logout)
	mux.HandleFunc("GET /auth/me", middleware.RequireAuth(h.jwtSecret)(h.Me))
	mux.HandleFunc("PUT /auth/password", authLimiter(middleware.RequireAuth(h.jwtSecret)(h.ChangePassword)))
	// The frontend's login page uses these now instead of POST /auth/login — that
	// route stays mounted for API back-compat (see AuthService.Login's callers).
	mux.HandleFunc("POST /auth/otp/request", authLimiter(h.RequestOTP))
	mux.HandleFunc("POST /auth/otp/verify", authLimiter(h.VerifyOTP))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, token, err := h.svc.Register(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	h.setSessionCookie(w, token)
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in dto.LoginRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, token, err := h.svc.Login(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	h.setSessionCookie(w, token)
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cookieSecure,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	user, err := h.svc.Me(r.Context(), authUser.ID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var in dto.ChangePasswordRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.ChangePassword(r.Context(), authUser.ID, in); err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var in dto.RequestOTPRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.RequestLoginOTP(r.Context(), in.Phone); err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var in dto.VerifyOTPRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, token, err := h.svc.VerifyLoginOTP(r.Context(), in.Phone, in.Code)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	h.setSessionCookie(w, token)
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cookieSecure,
		MaxAge:   int(auth.TTL.Seconds()),
	})
}
