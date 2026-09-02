package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/auth"
)

type contextKey string

const userContextKey contextKey = "authUser"

const CookieName = "session_token"

type AuthUser struct {
	ID   string
	Role string
}

// RequireAuth rejects the request with 401 unless it carries a valid session
// cookie, otherwise injects the authenticated user into the request context.
func RequireAuth(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := authenticate(r, secret)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "authentication required")
				return
			}
			next(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
		}
	}
}

// RequireRole is RequireAuth plus a 403 if the authenticated user's role doesn't match.
func RequireRole(secret, role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := authenticate(r, secret)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if user.Role != role {
				writeAuthError(w, http.StatusForbidden, "you don't have permission to do that")
				return
			}
			next(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
		}
	}
}

func authenticate(r *http.Request, secret string) (*AuthUser, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil, err
	}

	claims, err := auth.ParseToken(cookie.Value, secret)
	if err != nil {
		return nil, err
	}

	return &AuthUser{ID: claims.UserID, Role: claims.Role}, nil
}

func UserFromContext(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(*AuthUser)
	return user, ok
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
