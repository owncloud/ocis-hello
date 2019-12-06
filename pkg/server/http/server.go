package http

import (
	"github.com/owncloud/ocis-hello/pkg/assets"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/flagset"
	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis-pkg/middleware"
	"github.com/owncloud/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace("go.micro.web"),
		http.Name("hello"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(flagset.RootWithConfig(config.New())...),
		http.Flags(flagset.ServerWithConfig(config.New())...),
	)

	hello := svc.NewService()

	{
		hello = svc.NewInstrument(hello, options.Metrics)
		hello = svc.NewLogging(hello, options.Logger)
		hello = svc.NewTracing(hello)
	}

	proto.RegisterHelloWeb(
		service,
		hello,
		middleware.RealIP,
		middleware.RequestID,
		middleware.Cache,
		middleware.Cors,
		middleware.Secure,
		middleware.Version(
			"hello",
			version.String,
		),
		middleware.Logger(
			options.Logger,
		),
		middleware.Static(
			assets.New(
				assets.Logger(options.Logger),
				assets.Config(options.Config),
			),
		),
	)

	service.Init()
	return service, nil
}
