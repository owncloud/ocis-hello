package debug

import (
	"io"
	"net/http"

	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/service/debug"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Name),
		debug.Version(version.String),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(health(options.Logger)),
		debug.Ready(ready(options.Logger)),
	), nil
}

// health implements the health check.
func health(logger log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// TODO(tboerger): check if services are up and running

		if _, err := io.WriteString(w, http.StatusText(http.StatusOK)); err != nil {
			logger.Error().
				Err(err).
				Str("request", "health").
				Msg("Failed to write response")
		}
	}
}

// ready implements the ready check.
func ready(logger log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// TODO(tboerger): check if services are up and running

		if _, err := io.WriteString(w, http.StatusText(http.StatusOK)); err != nil {
			logger.Error().
				Err(err).
				Str("request", "ready").
				Msg("Failed to write response")
		}
	}
}
