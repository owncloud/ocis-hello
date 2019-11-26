package command

import (
	"fmt"
	"net/http"

	"github.com/micro/cli"
	"github.com/owncloud/ocis-hello/pkg/config"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "health",
		Usage: "Check health status",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "debug-addr",
				Value:       "0.0.0.0:8390",
				Usage:       "Address to debug endpoint",
				EnvVar:      "HELLO_DEBUG_ADDR",
				Destination: &cfg.Debug.Addr,
			},
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					cfg.Debug.Addr,
				),
			)

			if err != nil {
				logger.Fatal().
					Err(err).
					Msg("Failed to request health check")
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				logger.Fatal().
					Int("code", resp.StatusCode).
					Msg("Health seems to be in bad state")
			}

			logger.Debug().
				Int("code", resp.StatusCode).
				Msg("Health got a good state")

			return nil
		},
	}
}
