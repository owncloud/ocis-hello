package main

func main() {

}

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/go-kit/kit/log/level"
// 	"github.com/owncloud/ocis-hello/pkg/api/v0"
// 	"github.com/owncloud/ocis-hello/pkg/command"
// 	"github.com/owncloud/ocis-hello/pkg/config"
// 	"github.com/spf13/cobra"
// )

// var (
// 	ErrMissingName     = errors.New("missing a name to greet")
// 	ErrBadResponse     = errors.New("got an error as response")
// 	ErrUnknownProtocol = errors.New("unknown client protocol")
// )

// // Greet is the entrypoint for the greet command.
// func Greet(cfg *config.Config) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "greet <name>",
// 		Short: "Call greet API endpoint",
// 		Long:  ``,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			logger := command.NewLogger(cfg)

// 			if len(args) < 1 {
// 				level.Error(logger).Log(
// 					"msg", ErrMissingName,
// 				)

// 				return ErrMissingName
// 			}

// 			switch cfg.Client.Protocol {
// 			case "grpc":
// 				conn, err := GrpcClient(cfg)

// 				if err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to initialize client",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				defer conn.Close()

// 				client := v0api.NewHelloServiceClient(conn)

// 				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 				defer cancel()

// 				req := v0api.GreetRequest{
// 					Name: strings.Join(args, " "),
// 				}

// 				res, err := client.Greet(ctx, &req)

// 				if err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to submit request",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				level.Info(logger).Log(
// 					"msg", res.Message,
// 				)

// 			case "http":
// 				if !strings.HasSuffix(cfg.Client.Endpoint, "/") {
// 					cfg.Client.Endpoint = cfg.Client.Endpoint + "/"
// 				}

// 				client, err := RestClient(cfg)

// 				if err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to initialize client",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				content := struct {
// 					Name string `json:"name"`
// 				}{
// 					Name: strings.Join(args, " "),
// 				}

// 				buf := bytes.NewBufferString("")

// 				if err := json.NewEncoder(buf).Encode(content); err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to encode request",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				req, err := http.NewRequest(
// 					"POST",
// 					fmt.Sprintf(
// 						"%sapi/v0alpha1/greet",
// 						cfg.Client.Endpoint,
// 					),
// 					buf,
// 				)

// 				if err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to build request",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				resp, err := client.Do(req)

// 				if err != nil {
// 					level.Info(logger).Log(
// 						"msg", "Failed to submit request",
// 						"err", err,
// 					)

// 					return err
// 				}

// 				if resp.StatusCode != http.StatusOK {
// 					level.Error(logger).Log(
// 						"msg", http.StatusText(resp.StatusCode),
// 						"code", resp.StatusCode,
// 					)

// 					return ErrBadResponse
// 				} else {
// 					res := struct {
// 						Message string `json:"message"`
// 					}{}

// 					body, err := ioutil.ReadAll(resp.Body)

// 					if err != nil {
// 						level.Info(logger).Log(
// 							"msg", "Failed to read response",
// 							"err", err,
// 						)
// 					}

// 					if err := json.Unmarshal(body, &res); err != nil {
// 						level.Info(logger).Log(
// 							"msg", "Failed to parse response",
// 							"err", err,
// 						)
// 					}

// 					level.Info(logger).Log(
// 						"msg", res.Message,
// 					)
// 				}
// 			default:
// 				level.Error(logger).Log(
// 					"msg", "Unknown client protocol",
// 					"protocol", cfg.Client.Protocol,
// 				)

// 				return ErrUnknownProtocol
// 			}

// 			return nil
// 		},
// 	}

// 	cmd.Flags().String("client-protocol", "", "Protocol for API to use")
// 	cfg.Viper.BindPFlag("client.protocol", cmd.Flags().Lookup("client-protocol"))
// 	cfg.Viper.SetDefault("client.protocol", "grpc")

// 	cmd.Flags().String("client-endpoint", "", "Address to client endpoint")
// 	cfg.Viper.BindPFlag("client.endpoint", cmd.Flags().Lookup("client-endpoint"))
// 	cfg.Viper.SetDefault("client.endpoint", "localhost:8381")

// 	return cmd
// }
