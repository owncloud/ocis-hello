package command

import (
	"os"
	"os/signal"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/transport/debug"
	"github.com/owncloud/ocis-hello/pkg/transport/grpc"
	"github.com/owncloud/ocis-hello/pkg/transport/http"
	"github.com/spf13/cobra"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start integrated server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
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
						level.Error(logger).Log(
							"msg", "Failed to create agent tracing",
							"err", err,
							"endpoint", cfg.Tracing.Endpoint,
							"collector", cfg.Tracing.Collector,
						)

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
						level.Error(logger).Log(
							"msg", "Failed to create jaeger tracing",
							"err", err,
							"endpoint", cfg.Tracing.Endpoint,
							"collector", cfg.Tracing.Collector,
						)

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						level.Error(logger).Log(
							"msg", "Failed to create zipkin tracing",
							"err", err,
							"endpoint", cfg.Tracing.Endpoint,
							"collector", cfg.Tracing.Collector,
						)

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
					level.Warn(logger).Log(
						"msg", "Unknown tracing backend",
						"type", t,
					)
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				level.Debug(logger).Log(
					"msg", "Tracing is not enabled",
				)
			}

			var (
				gr run.Group
			)

			{
				server, err := debug.NewServer(
					debug.WithLogger(log.With(logger, "transport", "debug")),
					debug.WithAddr(cfg.Debug.Addr),
					debug.WithToken(cfg.Debug.Token),
					debug.WithPprof(cfg.Debug.Pprof),
				)

				if err != nil {
					level.Error(log.With(logger, "transport", "debug")).Log(
						"msg", "Failed to initialize server",
						"err", err,
					)

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(reason error) {
					server.Shutdown(reason)
				})
			}

			{
				server, err := http.NewServer(
					http.WithLogger(log.With(logger, "transport", "http")),
					http.WithAddr(cfg.HTTP.Addr),
					http.WithGrpc(cfg.GRPC.Addr),
					http.WithRoot(cfg.HTTP.Root),
					http.WithAssets(cfg.Asset.Path),
				)

				if err != nil {
					level.Error(log.With(logger, "transport", "http")).Log(
						"msg", "Failed to initialize server",
						"err", err,
					)

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(reason error) {
					server.Shutdown(reason)
				})
			}

			{
				server, err := grpc.NewServer(
					grpc.WithLogger(log.With(logger, "transport", "grpc")),
					grpc.WithAddr(cfg.GRPC.Addr),
				)

				if err != nil {
					level.Error(log.With(logger, "transport", "grpc")).Log(
						"msg", "Failed to initialize server",
						"err", err,
					)

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(reason error) {
					server.Shutdown(reason)
				})
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

	cmd.Flags().String("tracing-enabled", "", "Enable sending of traces")
	cfg.Viper.BindPFlag("tracing.enabled", cmd.Flags().Lookup("tracing-enabled"))
	cfg.Viper.SetDefault("tracing.enabled", false)

	cmd.Flags().String("tracing-type", "", "Tracing backend type")
	cfg.Viper.BindPFlag("tracing.type", cmd.Flags().Lookup("tracing-type"))
	cfg.Viper.SetDefault("tracing.type", "jaeger")

	cmd.Flags().String("tracing-endpoint", "", "Endpoint for the agent")
	cfg.Viper.BindPFlag("tracing.endpoint", cmd.Flags().Lookup("tracing-endpoint"))
	cfg.Viper.SetDefault("tracing.endpoint", "")

	cmd.Flags().String("tracing-collector", "", "Endpoint for the collector")
	cfg.Viper.BindPFlag("tracing.collector", cmd.Flags().Lookup("tracing-collector"))
	cfg.Viper.SetDefault("tracing.collector", "")

	cmd.Flags().String("tracing-service", "", "Service name for tracing")
	cfg.Viper.BindPFlag("tracing.service", cmd.Flags().Lookup("tracing-service"))
	cfg.Viper.SetDefault("tracing.service", "hello")

	cmd.Flags().String("debug-addr", "", "Address to bind debug server")
	cfg.Viper.BindPFlag("debug.addr", cmd.Flags().Lookup("debug-addr"))
	cfg.Viper.SetDefault("debug.addr", "0.0.0.0:8390")

	cmd.Flags().String("debug-token", "", "Token to grant metrics access")
	cfg.Viper.BindPFlag("debug.token", cmd.Flags().Lookup("debug-token"))
	cfg.Viper.SetDefault("debug.token", "")

	cmd.Flags().Bool("debug-pprof", false, "Enable pprof debugging")
	cfg.Viper.BindPFlag("debug.pprof", cmd.Flags().Lookup("debug-pprof"))
	cfg.Viper.SetDefault("debug.pprof", false)

	cmd.Flags().String("http-addr", "", "Address to bind http server")
	cfg.Viper.BindPFlag("http.addr", cmd.Flags().Lookup("http-addr"))
	cfg.Viper.SetDefault("http.addr", "0.0.0.0:8380")

	cmd.Flags().String("http-root", "", "Root path for http endpoint")
	cfg.Viper.BindPFlag("http.root", cmd.Flags().Lookup("http-root"))
	cfg.Viper.SetDefault("http.root", "/")

	cmd.Flags().String("grpc-addr", "", "Address to bind grpc server")
	cfg.Viper.BindPFlag("grpc.addr", cmd.Flags().Lookup("grpc-addr"))
	cfg.Viper.SetDefault("grpc.addr", "0.0.0.0:8381")

	cmd.Flags().String("grpc-root", "", "Root path for grpc endpoint")
	cfg.Viper.BindPFlag("grpc.root", cmd.Flags().Lookup("grpc-root"))
	cfg.Viper.SetDefault("grpc.root", "/")

	cmd.Flags().String("asset-path", "", "Path to custom assets")
	cfg.Viper.BindPFlag("asset.path", cmd.Flags().Lookup("asset-path"))
	cfg.Viper.SetDefault("asset.path", "")

	return cmd
}
