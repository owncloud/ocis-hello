package grpc

import (
	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis-pkg/service/grpc"
)

// Server initializes the grpc service and server.
func Server(opts ...Option) (grpc.Service, error) {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Namespace(options.Namespace),
		grpc.Name("api.hello"),
		grpc.Version(version.String),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
	)

	var hello proto.HelloHandler
	{
		hello = svc.NewService()
		hello = svc.NewInstrument(hello, options.Metrics)
		hello = svc.NewLogging(hello, options.Logger)
	}

	proto.RegisterHelloHandler(
		service.Server(),
		hello,
	)

	service.Init()
	return service, nil
}
