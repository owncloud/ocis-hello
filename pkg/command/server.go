package command

import (
	"context"
	"errors"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/flagset"
	"github.com/owncloud/ocis-hello/pkg/metrics"
	"github.com/owncloud/ocis-hello/pkg/server/grpc"
	"github.com/owncloud/ocis-hello/pkg/server/http"
	"github.com/owncloud/ocis-hello/pkg/service/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/sync"
	smessages "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settings "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	ssvc "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	client "go-micro.dev/v4/client"
	"strings"
	"time"
)

const (
	bundleIDGreeting       = "21fb587b-7b69-4aa6-b0a7-93c74af1918f"
	settingIDGreeterPhrase = "b3584ea8-caec-4951-a2c1-92cbc70071b7"

	// maxRetries indicates how many times to try a request for network reasons.
	maxRetries = 5
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "start hello service",
		Description: "Hello is an example oCIS extension",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			logger := NewLogger(cfg)
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			if !cfg.Supervised {
				return ParseConfig(ctx, cfg)
			}
			logger.Debug().Str("service", "hello").Msg("ignoring config file parsing when running supervised")
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			err := tracing.Configure(cfg, logger)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)
			mtrcs := metrics.New()

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			bundleService := settings.NewBundleService("com.owncloud.api.settings", client.DefaultClient)

			for i := 1; i <= maxRetries; i++ {
				err = registerSettingsBundles(bundleService, &logger)
				if err != nil {
					logger.Logger.Info().Msg(err.Error())
					// limited potential backoff: 1s, 4s, 9s, 16s, 25s, ..., but max 30s
					backoff := time.Duration(i*i) * time.Second
					if backoff > 30*time.Second {
						backoff = 30 * time.Second
					}
					logger.Logger.Info().Dur("backoff", backoff).Msg("retry to register settings bundle and permission")
					time.Sleep(backoff)
				}
			}
			if err != nil {
				logger.Error().Msg("failed to register settings - aborting server initialization")
				logger.Info().Msg("shutting down server")
				cancel()
			}

			ps := settingsPhraseSource{vsClient: settings.NewValueService("com.owncloud.api.settings", client.DefaultClient)}
			handler, err := svc.NewGreeter(svc.PhraseSource(ps), svc.Logger(logger))
			if err != nil {
				logger.Error().Err(err).Msg("handler init")
				return err
			}

			handler = service.NewInstrument(handler, mtrcs)
			handler = service.NewLogging(handler, logger)
			handler = service.NewTracing(handler)

			httpServer := http.Server(
				http.Config(cfg),
				http.Logger(logger),
				http.Name(cfg.Server.Name),
				http.Context(ctx),
				http.Metrics(mtrcs),
				http.Handler(handler),
			)

			gr.Add(httpServer.Run, func(_ error) {
				logger.Info().Str("server", "http").Msg("shutting down server")
				cancel()
			})

			grpcServer := grpc.Server(
				grpc.Config(cfg),
				grpc.Logger(logger),
				grpc.Name(cfg.Server.Name),
				grpc.Context(ctx),
				grpc.Metrics(mtrcs),
				grpc.Handler(handler),
			)

			gr.Add(grpcServer.Run, func(_ error) {
				logger.Info().Str("server", "grpc").Msg("shutting down server")
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// defineContext sets the context for the extension. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}

// registerSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func registerSettingsBundles(bundleService settings.BundleService, l *log.Logger) (err error) {

	request := &settings.SaveBundleRequest{
		Bundle: &smessages.Bundle{
			Id:          bundleIDGreeting,
			Name:        "greeting",
			DisplayName: "Greeting",
			Extension:   "ocis-hello",
			Type:        smessages.Bundle_TYPE_DEFAULT,
			Resource: &smessages.Resource{
				Type: smessages.Resource_TYPE_SYSTEM,
			},
			Settings: []*smessages.Setting{
				{
					Id:          settingIDGreeterPhrase,
					Name:        "phrase",
					DisplayName: "Phrase",
					Description: "Phrase for replies on the greet request",
					Resource: &smessages.Resource{
						Type: smessages.Resource_TYPE_SYSTEM,
					},
					Value: &smessages.Setting_StringValue{
						StringValue: &smessages.String{
							Required:  true,
							Default:   "Hello",
							MaxLength: 15,
						},
					},
				},
			},
		},
	}

	_, err = bundleService.SaveBundle(context.Background(), request)
	if err != nil {
		l.Logger.Info().Msg("11111-------" + err.Error())
		l.With().Err(err).Logger().With().Str("settings bundle ID", request.Bundle.Id).Err(errors.New("could not create / update the settings bundle"))
		return err
	}
	l.Logger.Info().Str("settings bundle ID", request.Bundle.Id).Msg("created / updated the settings bundle")

	permissionRequests := []*settings.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &smessages.Setting{
				Id: "d5f42c4b-e1b6-4b59-8eca-fc4b9e9f2320",
				Resource: &smessages.Resource{
					Type: smessages.Resource_TYPE_SETTING,
					Id:   settingIDGreeterPhrase,
				},
				Name: "phrase-admin-read",
				Value: &smessages.Setting_PermissionValue{
					PermissionValue: &smessages.Permission{
						Operation:  smessages.Permission_OPERATION_READWRITE,
						Constraint: smessages.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}

	for _, permissionRequest := range permissionRequests {
		_, err := bundleService.AddSettingToBundle(context.Background(), permissionRequest)
		if err != nil {
			l.Logger.Info().Msg(err.Error())
			l.With().Err(err).Logger().With().Str("permission request bundle ID", request.Bundle.Id).Err(errors.New("could not create / update the permissions of the settings bundle"))
			return err
		}
		l.Logger.Info().Str("permission request bundle ID", request.Bundle.Id).Msg("created / updated the permissions of the settings bundle")
	}

	return nil
}

type settingsPhraseSource struct {
	vsClient settings.ValueService
}

func (s settingsPhraseSource) GetPhrase(accountID string) string {
	// request to the settings service requires to have the account uuid of the authenticated user available in the context
	rq := settings.GetValueByUniqueIdentifiersRequest{
		AccountUuid: accountID,
		SettingId:   settingIDGreeterPhrase,
	}

	response, err := s.vsClient.GetValueByUniqueIdentifiers(context.Background(), &rq)
	if err == nil {
		value, ok := response.Value.Value.Value.(*smessages.Value_StringValue)
		if ok {
			trimmedPhrase := strings.Trim(
				value.StringValue,
				" \t",
			)
			if trimmedPhrase != "" {
				return trimmedPhrase + " %s"
			}
		}
	}
	return service.DefaultPhrase

}
