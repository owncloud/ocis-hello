package command

import (
	"os"
	"strings"

	"github.com/micro/cli"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-hello command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis-hello",
		Version:  version.String,
		Usage:    "Example service for Reva/oCIS",
		Compiled: version.Compiled(),

		Authors: []cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Value:       "",
				Usage:       "Path to config file",
				EnvVar:      "HELLO_CONFIG_FILE",
				Destination: &cfg.File,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Usage:       "Set logging level",
				EnvVar:      "HELLO_LOG_LEVEL",
				Destination: &cfg.Log.Level,
			},
			&cli.BoolTFlag{
				Name:        "log-pretty",
				Usage:       "Enable pretty logging",
				EnvVar:      "HELLO_LOG_PRETTY",
				Destination: &cfg.Log.Pretty,
			},
			&cli.BoolTFlag{
				Name:        "log-color",
				Usage:       "Enable colored logging",
				EnvVar:      "HELLO_LOG_COLOR",
				Destination: &cfg.Log.Color,
			},
		},

		Before: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetEnvPrefix("HELLO")
			viper.AutomaticEnv()

			if c.IsSet("config-file") {
				viper.SetConfigFile(c.String("config-file"))
			} else {
				viper.SetConfigName("hello")

				viper.AddConfigPath("/etc/ocis")
				viper.AddConfigPath("$HOME/.ocis")
				viper.AddConfigPath("./config")
			}

			if err := viper.ReadInConfig(); err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					logger.Info().
						Msg("Continue without config")
				case viper.UnsupportedConfigError:
					logger.Fatal().
						Err(err).
						Msg("Unsupported config type")

					return err
				default:
					logger.Fatal().
						Err(err).
						Msg("Failed to read config")

					return err
				}
			}

			if err := viper.Unmarshal(&cfg); err != nil {
				logger.Fatal().
					Err(err).
					Msg("Failed to parse config")
			}

			return nil
		},

		Commands: []cli.Command{
			Server(cfg),
			Health(cfg),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version,v",
		Usage: "Print the version",
	}

	return app.Run(os.Args)
}

func NewLogger(cfg *config.Config) zerolog.Logger {
	switch strings.ToLower(cfg.Log.Level) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if cfg.Log.Pretty {
		log.Logger = log.Output(
			zerolog.NewConsoleWriter(
				func(w *zerolog.ConsoleWriter) {
					w.Out = os.Stderr
					w.NoColor = !cfg.Log.Color
				},
			),
		)
	}

	return log.Logger.With().Str("service", "hello").Timestamp().Logger()
}
