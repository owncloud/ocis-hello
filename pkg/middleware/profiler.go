package middleware

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

// Profiler is a convenient subrouter used for mounting net/http/pprof.
func Profiler() http.Handler {
	return middleware.Profiler()
}
