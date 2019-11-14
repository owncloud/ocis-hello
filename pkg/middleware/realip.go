package middleware

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

// RealIP is a middleware that sets a http.Request RemoteAddr.
func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}
