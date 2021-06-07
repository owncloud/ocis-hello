package grpc

import (
	"context"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)
	greeter := options.Handler

	service := grpc.NewService(
		grpc.Name(options.Config.Server.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
		grpc.Version(options.Config.Server.Version),
	)

	if err := proto.RegisterHelloHandler(service.Server(), helloHandler{greeter: greeter}); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	service.Init()
	return service
}

type helloHandler struct {
	greeter service.Greeter
}

func (h helloHandler) Greet(ctx context.Context, req *proto.GreetRequest, rsp *proto.GreetResponse) error {
	name := req.Name
	if name == "" {
		return service.ErrMissingName
	}

	// TODO: can we get the current users accountID somehow?
	rsp.Message = h.greeter.Greet("", name)

	return nil
}
