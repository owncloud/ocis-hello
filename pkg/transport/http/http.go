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
	"github.com/justinas/alice"
	"github.com/owncloud/ocis-hello/pkg/api/v0alpha1"
	"github.com/owncloud/ocis-hello/pkg/middleware"
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
		Addr: s.addr,
		Handler: alice.New(
			middleware.RequestID,
			middleware.RealIP,
			middleware.Version,
			middleware.Cache,
			middleware.Secure,
			middleware.Options,

			middleware.Logger(
				middleware.WithLogger(s.logger),
			),
			middleware.Static(
				middleware.WithLogger(s.logger),
				middleware.WithRoot(s.root),
				middleware.WithAssets(s.assets),
			),
		).Then(mux),
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
