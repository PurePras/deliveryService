package middleware

import "net/http"

// SecurityHeaders sets a small set of headers relevant to a JSON API. Content-Security-Policy
// is deliberately not set here — this backend never serves HTML, so CSP belongs to the
// frontend's eventual nginx config instead.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
