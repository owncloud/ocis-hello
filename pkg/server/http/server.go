package http

import (
	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-hello/pkg/assets"
	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis-pkg/middleware"
	"github.com/owncloud/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace("com.owncloud.web"),
		http.Name("hello"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	hello := svc.NewService()

	{
		hello = svc.NewInstrument(hello, options.Metrics)
		hello = svc.NewLogging(hello, options.Logger)
		hello = svc.NewTracing(hello)
	}

	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Cache)
	mux.Use(middleware.Cors)
	mux.Use(middleware.Secure)

	mux.Use(middleware.Version(
		"hello",
		version.String,
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Use(middleware.Static(
		options.Config.HTTP.Root,
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		proto.RegisterHelloWeb(
			r,
			hello,
		)
	})

	service.Handle(
		"/",
		mux,
	)

	service.Init()
	return service, nil
}
