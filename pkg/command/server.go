package command

import (
	"os"
	"os/signal"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis-hello/pkg/config"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	// 	"github.com/owncloud/ocis-hello/pkg/config"
	// 	"github.com/owncloud/ocis-hello/pkg/transport/debug"
	// 	"github.com/owncloud/ocis-hello/pkg/transport/grpc"
	"github.com/owncloud/ocis-hello/pkg/server/http"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "tracing-enabled",
				Usage:       "Enable sending traces",
				EnvVar:      "HELLO_TRACING_ENABLED",
				Destination: &cfg.Tracing.Enabled,
			},
			cli.StringFlag{
				Name:        "tracing-type",
				Value:       "jaeger",
				Usage:       "Tracing backend type",
				EnvVar:      "HELLO_TRACING_TYPE",
				Destination: &cfg.Tracing.Type,
			},
			cli.StringFlag{
				Name:        "tracing-endpoint",
				Value:       "",
				Usage:       "Endpoint for the agent",
				EnvVar:      "HELLO_TRACING_ENDPOINT",
				Destination: &cfg.Tracing.Endpoint,
			},
			cli.StringFlag{
				Name:        "tracing-collector",
				Value:       "",
				Usage:       "Endpoint for the collector",
				EnvVar:      "HELLO_TRACING_COLLECTOR",
				Destination: &cfg.Tracing.Collector,
			},
			cli.StringFlag{
				Name:        "tracing-service",
				Value:       "hello",
				Usage:       "Service name for tracing",
				EnvVar:      "HELLO_TRACING_SERVICE",
				Destination: &cfg.Tracing.Service,
			},
			cli.StringFlag{
				Name:        "debug-addr",
				Value:       "0.0.0.0:8390",
				Usage:       "Address to bind debug server",
				EnvVar:      "HELLO_DEBUG_ADDR",
				Destination: &cfg.Debug.Addr,
			},
			cli.StringFlag{
				Name:        "debug-token",
				Value:       "",
				Usage:       "Token to grant metrics access",
				EnvVar:      "HELLO_DEBUG_TOKEN",
				Destination: &cfg.Debug.Token,
			},
			cli.BoolFlag{
				Name:        "debug-pprof",
				Usage:       "Enable pprof debugging",
				EnvVar:      "HELLO_DEBUG_PPROF",
				Destination: &cfg.Debug.Pprof,
			},
			cli.StringFlag{
				Name:        "http-addr",
				Value:       "0.0.0.0:8380",
				Usage:       "Address to bind http server",
				EnvVar:      "HELLO_HTTP_ADDR",
				Destination: &cfg.HTTP.Addr,
			},
			cli.StringFlag{
				Name:        "grpc-addr",
				Value:       "0.0.0.0:8381",
				Usage:       "Address to bind grpc server",
				EnvVar:      "HELLO_GRPC_ADDR",
				Destination: &cfg.GRPC.Addr,
			},
			cli.StringFlag{
				Name:        "asset-path",
				Value:       "",
				Usage:       "Path to custom assets",
				EnvVar:      "HELLO_ASSET_PATH",
				Destination: &cfg.Asset.Path,
			},
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					exporter, err := ocagent.NewExporter(
						ocagent.WithReconnectionPeriod(5*time.Second),
						ocagent.WithAddress(cfg.Tracing.Endpoint),
						ocagent.WithServiceName(cfg.Tracing.Service),
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create agent tracing")

						return err
					}

					trace.RegisterExporter(exporter)
					view.RegisterExporter(exporter)

				case "jaeger":
					exporter, err := jaeger.NewExporter(
						jaeger.Options{
							AgentEndpoint:     cfg.Tracing.Endpoint,
							CollectorEndpoint: cfg.Tracing.Collector,
							ServiceName:       cfg.Tracing.Service,
						},
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create jaeger tracing")

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create zipkin tracing")

						return err
					}

					exporter := zipkin.NewExporter(
						zipkinhttp.NewReporter(
							cfg.Tracing.Collector,
						),
						endpoint,
					)

					trace.RegisterExporter(exporter)

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var gr run.Group

			{
				// server, err := debug.NewServer(
				// 	debug.WithLogger(log.With(logger, "transport", "debug")),
				// 	debug.WithAddr(cfg.Debug.Addr),
				// 	debug.WithToken(cfg.Debug.Token),
				// 	debug.WithPprof(cfg.Debug.Pprof),
				// )

				// if err != nil {
				// 	level.Error(log.With(logger, "transport", "debug")).Log(
				// 		"msg", "Failed to initialize server",
				// 		"err", err,
				// 	)

				// 	return err
				// }

				// gr.Add(func() error {
				// 	return server.ListenAndServe()
				// }, func(reason error) {
				// 	server.Shutdown(reason)
				// })
			}

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Addr(cfg.HTTP.Addr),
					http.Assets(cfg.Asset.Path),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					if err := server.Run(); err != nil {
						logger.Info().
							Err(err).
							Str("transport", "http").
							Msg("Failed to start server")

						return err
					}

					return nil
				}, func(reason error) {
					logger.Info().
						Str("transport", "http").
						Msg("Shutting down server")
				})
			}

			{
				// server, err := grpc.NewServer(
				// 	grpc.WithLogger(log.With(logger, "transport", "grpc")),
				// 	grpc.WithAddr(cfg.GRPC.Addr),
				// )

				// if err != nil {
				// 	level.Error(log.With(logger, "transport", "grpc")).Log(
				// 		"msg", "Failed to initialize server",
				// 		"err", err,
				// 	)

				// 	return err
				// }

				// gr.Add(func() error {
				// 	return server.ListenAndServe()
				// }, func(reason error) {
				// 	server.Shutdown(reason)
				// })
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
				})
			}

			return gr.Run()
		},
	}
}
