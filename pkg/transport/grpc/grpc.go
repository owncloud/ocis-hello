package grpc

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/owncloud/ocis-hello/pkg/api/v0alpha1"
	"github.com/owncloud/ocis-hello/pkg/endpoint/v0alpha1"
	"github.com/owncloud/ocis-hello/pkg/service/v0alpha1"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type Server struct {
	logger log.Logger
	addr   string

	server *grpc.Server
}

func (s *Server) ListenAndServe() error {
	level.Info(s.logger).Log(
		"msg", "Starting server",
		"addr", s.addr,
	)

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			kitgrpc.Interceptor,
		),
	)

	counter := prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "hello",
		Subsystem: "greet",
		Name:      "counter",
		Help:      "Total requests for a greeting.",
	}, []string{"method", "success"})

	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "hello",
		Subsystem: "greet",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds.",
	}, []string{"method", "success"})

	service := v0alpha1svc.New(
		v0alpha1svc.WithLogger(s.logger),
		v0alpha1svc.WithCounter(counter),
	)

	endpoint := v0alpha1endpoint.New(
		v0alpha1endpoint.WithLogger(s.logger),
		v0alpha1endpoint.WithService(service),
		v0alpha1endpoint.WithDuration(duration),
	)

	v0alpha1api.RegisterHelloServiceServer(
		s.server,
		NewV0alpha1(
			endpoint,
			s.logger,
		),
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

	listener, err := net.Listen(
		"tcp",
		s.addr,
	)

	if err != nil {
		level.Error(s.logger).Log(
			"msg", "Failed to initialize listener",
			"err", err,
			"addr", s.addr,
		)

		return err
	}

	return s.server.Serve(listener)
}

func (s *Server) Shutdown(reason error) {
	s.server.GracefulStop()

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

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
