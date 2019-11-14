package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/owncloud/ocis-hello/pkg/api/v0alpha1"
	// "github.com/owncloud/ocis-hello/pkg/handler/static"
	// "github.com/owncloud/ocis-hello/pkg/middleware"
	"google.golang.org/grpc"
)

type Server struct {
	logger log.Logger
	addr   string
	grpc   string
	root   string
	assets string

	server *http.Server
}

func (s *Server) ListenAndServe() error {
	level.Info(s.logger).Log(
		"msg", "Starting server",
		"addr", s.addr,
	)

	ctx := context.Background()
	mux := runtime.NewServeMux()

	if err := v0alpha1api.RegisterHelloServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpc,
		[]grpc.DialOption{
			grpc.WithInsecure(),
		},
	); err != nil {
		level.Error(s.logger).Log(
			"msg", "Failed to bind gRPC server",
			"api", "v0alpha1",
			"err", err,
		)

		return err
	}

	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

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

func (s *Server) Shutdown(reason error) {
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

func NewServer(opts ...Option) (*Server, error) {
	s := new(Server)

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// func (s Server) Router() *chi.Mux {
// 	mux := chi.NewRouter()

// 	// mux.Use(hlog.NewHandler(log.Logger))
// 	// mux.Use(hlog.RemoteAddrHandler("ip"))
// 	// mux.Use(hlog.URLHandler("path"))
// 	// mux.Use(hlog.MethodHandler("method"))
// 	// mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

// 	// mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
// 	// 	hlog.FromRequest(r).Debug().
// 	// 		Str("method", r.Method).
// 	// 		Str("url", r.URL.String()).
// 	// 		Int("status", status).
// 	// 		Int("size", size).
// 	// 		Dur("duration", duration).
// 	// 		Msg("")
// 	// }))

// 	mux.Use(middleware.RealIP)
// 	mux.Use(middleware.Version)
// 	mux.Use(middleware.Cache)
// 	mux.Use(middleware.Secure)
// 	mux.Use(middleware.Options)

// 	mux.Route(s.root, func(root chi.Router) {
// 		root.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Content-Type", "text/plain")
// 			w.WriteHeader(http.StatusNotFound)

// 			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		})

// 		root.Route("/api", func(base chi.Router) {
// 			// base.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
// 			// 	request := struct {
// 			// 		Name string `json:"name"`
// 			// 	}{}

// 			// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 			// 		log.Error().
// 			// 			Err(err).
// 			// 			Msg("Failed to parse request")

// 			// 		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
// 			// 		return
// 			// 	}

// 			// 	payload := struct {
// 			// 		Message string `json:"message"`
// 			// 	}{}

// 			// 	payload.Message = fmt.Sprintf("Hello %s", request.Name)
// 			// 	response, err := json.Marshal(payload)

// 			// 	if err != nil {
// 			// 		log.Error().
// 			// 			Err(err).
// 			// 			Msg("Failed to build response")

// 			// 		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
// 			// 		return
// 			// 	}

// 			// 	w.Header().Set("Content-Type", "application/json")
// 			// 	w.WriteHeader(http.StatusOK)
// 			// 	w.Write(response)
// 			// })
// 		})

// 		root.Mount(
// 			"/",
// 			static.Handler(
// 				static.WithRoot(s.root),
// 				static.WithPath(s.assets),
// 			),
// 		)
// 	})

// 	return mux
// }
