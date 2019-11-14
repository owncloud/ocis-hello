package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/spf13/cobra"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check health status",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			logger := NewLogger(cfg)

			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					cfg.Debug.Addr,
				),
			)

			if err != nil {
				level.Error(logger).Log(
					"msg", "Failed to request health check",
					"err", err,
				)

				os.Exit(1)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				level.Error(logger).Log(
					"msg", "Health seems to be in bad state",
					"code", resp.StatusCode,
				)

				os.Exit(1)
			}

			level.Debug(logger).Log(
				"msg", "Health got a good state",
				"code", resp.StatusCode,
			)

			os.Exit(0)
		},
	}

	cmd.Flags().String("debug-addr", "", "Address to debug endpoint")
	cfg.Viper.BindPFlag("debug.addr", cmd.Flags().Lookup("debug-addr"))
	cfg.Viper.SetDefault("debug.addr", "0.0.0.0:8390")

	return cmd
}
