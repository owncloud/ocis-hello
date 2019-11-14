package debug

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/owncloud/ocis-hello/pkg/middleware"
)

type Server struct {
	logger log.Logger
	addr   string
	token  string
	pprof  bool

	server  *http.Server
	metrics *prometheus.Exporter
}

func (s Server) ListenAndServe() error {
	level.Info(s.logger).Log(
		"msg", "Starting server",
		"addr", s.addr,
	)

	if strings.HasPrefix(s.addr, "unix://") {
		socket := strings.TrimPrefix(s.addr, "unix://")

		if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
			level.Error(s.logger).Log(
				"msg", "Failed to remove existing socket",
				"err", err,
				"socket", socket,
			)

			return err
		}

		listener, err := net.ListenUnix(
			"unix",
			&net.UnixAddr{
				Name: socket,
				Net:  "unix",
			},
		)

		if err != nil {
			level.Error(s.logger).Log(
				"msg", "Failed to initialize unix socket",
				"err", err,
				"socket", socket,
			)

			return err
		}

		if err = os.Chmod(socket, os.FileMode(0666)); err != nil {
			level.Error(s.logger).Log(
				"msg", "Failed to change socket permissions",
				"err", err,
				"socket", socket,
			)

			return err
		}

		return s.server.Serve(listener)
	}

	return s.server.ListenAndServe()
}

func (s Server) Shutdown(reason error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		level.Error(s.logger).Log(
			"msg", "Failed to shutdown server gracefully",
			"err", err,
		)

		return
	}

	if strings.HasPrefix(s.addr, "unix://") {
		socket := strings.TrimPrefix(s.addr, "unix://")

		if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
			level.Error(s.logger).Log(
				"msg", "Failed to remove server socket",
				"err", err,
				"socket", socket,
			)
		}
	}

	level.Info(s.logger).Log(
		"msg", "Shutdown server gracefully",
		"reason", reason,
	)
}

func (s Server) Router() *chi.Mux {
	mux := chi.NewRouter()

	// mux.Use(hlog.NewHandler(log.Logger))
	// mux.Use(hlog.RemoteAddrHandler("ip"))
	// mux.Use(hlog.URLHandler("path"))
	// mux.Use(hlog.MethodHandler("method"))
	// mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(middleware.RealIP)
	mux.Use(middleware.Version)
	mux.Use(middleware.Cache)
	mux.Use(middleware.Secure)
	mux.Use(middleware.Options)

	mux.Route("/", func(root chi.Router) {
		if s.pprof {
			root.Mount(
				"/debug",
				middleware.Profiler(),
			)
		}

		root.With(
			middleware.Token(s.token),
		).Mount(
			"/metrics",
			s.metrics,
		)

		root.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}

func NewServer(opts ...Option) (*Server, error) {
	s := new(Server)

	for _, opt := range opts {
		opt(s)
	}

	if metrics, err := prometheus.NewExporter(
		prometheus.Options{
			Namespace: "hello",
		},
	); err != nil {
		return nil, err
	} else {
		s.metrics = metrics
	}

	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      s.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s, nil
}
