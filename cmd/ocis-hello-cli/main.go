package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/jackspirou/syscerts"
	"github.com/owncloud/ocis-hello/pkg/command"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func main() {
	if err := Root().Execute(); err != nil {
		os.Exit(1)
	}
}

// Root is the entry point for the ocis-hello-cli command.
func Root() *cobra.Command {
	cfg := config.New()

	cmd := &cobra.Command{
		Use:           "ocis-hello-cli",
		Short:         "Access the ocis-hello API",
		Long:          ``,
		Version:       version.String,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := command.NewLogger(cfg)

			if err := cfg.Viper.Unmarshal(&cfg); err != nil {
				level.Error(logger).Log(
					"msg", "Failed to process config",
					"err", err,
				)

				return err
			}

			return nil
		},
	}

	cfg.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.Viper.SetEnvPrefix("HELLO")
	cfg.Viper.AutomaticEnv()

	cmd.PersistentFlags().String("log-level", "", "Set logging level")
	cfg.Viper.BindPFlag("log.level", cmd.PersistentFlags().Lookup("log-level"))
	cfg.Viper.SetDefault("log.level", "info")

	cmd.PersistentFlags().Bool("log-pretty", false, "Enable pretty logging")
	cfg.Viper.BindPFlag("log.pretty", cmd.PersistentFlags().Lookup("log-pretty"))
	cfg.Viper.SetDefault("log.pretty", true)

	cmd.AddCommand(Greet(cfg))

	return cmd
}

func GrpcClient(cfg *config.Config) (*grpc.ClientConn, error) {
	return grpc.Dial(
		cfg.Client.Endpoint,
		grpc.WithInsecure(),
	)
}

func RestClient(cfg *config.Config) (*http.Client, error) {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs: syscerts.SystemRootsPool(),
			},
		},
	}, nil
}
