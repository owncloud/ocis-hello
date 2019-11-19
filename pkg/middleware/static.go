package middleware

import (
	"net/http"
	"strings"

	"github.com/owncloud/ocis-hello/pkg/assets"
)

// Static is a middleware that serves static assets.
func Static(opts ...Option) func(http.Handler) http.Handler {
	o := new(options)

	for _, opt := range opts {
		opt(o)
	}

	static := http.StripPrefix(
		o.root,
		http.FileServer(
			assets.New(
				assets.WithLogger(o.logger),
				assets.WithPath(o.assets),
			),
		),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				if strings.HasSuffix(r.URL.Path, "/") {
					http.NotFound(w, r)
				} else {
					static.ServeHTTP(w, r)
				}
			}
		})
	}
}
