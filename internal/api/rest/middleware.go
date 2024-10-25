package rest

import (
	"log/slog"
	"net/http"
	"time"
)

// LogMiddleware returns new request logging HTTPMiddleware.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rw := &responseWrapper{
			ResponseWriter: w,
		}

		start := time.Now()

		defer func() {
			status := rw.status

			values := []any{
				"http.method", r.Method,
				"http.url", r.URL.String(),
				"http.status", status,
				"http.latency", time.Since(start).String(),
			}

			if status >= http.StatusBadRequest && status <= http.StatusNetworkAuthenticationRequired {
				slog.ErrorContext(ctx, http.StatusText(status), values...)
			} else {
				slog.InfoContext(ctx, "", values...)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.status = code

	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	return rw.ResponseWriter.Write(b)
}

func (rw *responseWrapper) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}
