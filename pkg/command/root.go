package command

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root is the entry point for the ocis-hello command.
func Root() *cobra.Command {
	cfg := config.New()

	cmd := &cobra.Command{
		Use:          "ocis-hello",
		Short:        "Reva service for helloworld",
		Long:         ``,
		Version:      version.String,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(cfg)

			cfg.Viper.SetConfigName("hello")

			cfg.Viper.AddConfigPath("/etc/ocis")
			cfg.Viper.AddConfigPath("$HOME/.ocis")
			cfg.Viper.AddConfigPath("./config")

			if err := cfg.Viper.ReadInConfig(); err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					level.Debug(logger).Log(
						"msg", "Continue without config",
					)
				case viper.UnsupportedConfigError:
					level.Error(logger).Log(
						"msg", "Unsupported config type",
					)

					return err
				default:
					level.Error(logger).Log(
						"msg", "Failed to read config",
						"err", err,
					)

					return err
				}
			}

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
	cfg.Viper.SetDefault("log.pretty", false)

	cmd.AddCommand(Server(cfg))
	cmd.AddCommand(Health(cfg))

	return cmd
}

func NewLogger(cfg *config.Config) log.Logger {
	var (
		logger log.Logger
	)

	if cfg.Viper.GetBool("log.pretty") {
		logger = log.NewSyncLogger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)))
	} else {
		logger = log.NewSyncLogger(log.NewJSONLogger(log.NewSyncWriter(os.Stdout)))
	}

	switch strings.ToLower(cfg.Viper.GetString("log.level")) {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return log.With(
		logger,
		"time", log.DefaultTimestampUTC,
		"service", "hello",
	)
}
