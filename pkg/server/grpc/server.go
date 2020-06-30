package grpc

import (
	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
)

// Server initializes the grpc service and server.
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Name(options.Name),
		grpc.Version(version.String),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
	)

	handler := svc.NewService()
	handler = svc.NewInstrument(handler, options.Metrics)
	handler = svc.NewLogging(handler, options.Logger)
	if err := proto.RegisterHelloHandler(service.Server(), handler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Hello service handler")
	}

	service.Init()
	return service
}
