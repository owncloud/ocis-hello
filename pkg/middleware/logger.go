package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log/level"
)

// Logger is a middleware to log http requests.
func Logger(opts ...Option) func(http.Handler) http.Handler {
	o := new(options)

	for _, opt := range opts {
		opt(o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrap, r)

			level.Debug(o.logger).Log(
				"request", r.Header.Get("X-Request-ID"),
				"proto", r.Proto,
				"method", r.Method,
				"status", wrap.Status(),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", wrap.BytesWritten(),
			)
		})
	}
}
