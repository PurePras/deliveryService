package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder captures the status code a handler responds with. It has to override
// both WriteHeader and Write: most handlers call WriteHeader explicitly, but a handler
// that only calls Write (e.g. the health check) implicitly sends 200 on the first Write —
// on the *underlying* ResponseWriter, not through this wrapper — so Write needs its own
// fallback or that case would silently log status 0.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

// RequestLogger logs one structured line per request: method, path, status, duration, remote address.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_addr", clientIP(r),
			)
		})
	}
}
