package command

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	"github.com/owncloud/ocis/ocis-pkg/log"
	ogrpc "github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis/settings/pkg/service/v0"
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

			bundleService := settings.NewBundleService("com.owncloud.api.settings", ogrpc.DefaultClient)

			registerSettingsBundles(bundleService, &logger)

			ps := settingsPhraseSource{vsClient: settings.NewValueService("com.owncloud.api.settings", ogrpc.DefaultClient)}
			handler, err := svc.NewGreeter(svc.PhraseSource(ps), svc.Logger(logger))
			if err != nil {
				logger.Error().Err(err).Msg("handler init")
				return err
			}

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
func registerSettingsBundles(bundleService settings.BundleService, l *log.Logger) {
	request := &settings.SaveBundleRequest{
		Bundle: &settings.Bundle{
			Id:          bundleIDGreeting,
			Name:        "greeting",
			DisplayName: "Greeting",
			Extension:   "ocis-hello",
			Type:        settings.Bundle_TYPE_DEFAULT,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			Settings: []*settings.Setting{
				{
					Id:          settingIDGreeterPhrase,
					Name:        "phrase",
					DisplayName: "Phrase",
					Description: "Phrase for replies on the greet request",
					Resource: &settings.Resource{
						Type: settings.Resource_TYPE_SYSTEM,
					},
					Value: &settings.Setting_StringValue{
						StringValue: &settings.String{
							Required:  true,
							Default:   "Hello",
							MaxLength: 15,
						},
					},
				},
			},
		},
	}

	response, err := bundleService.SaveBundle(context.Background(), request)
	if err != nil {
		l.Warn().Msg("error registering settings bundle at first try. retrying")
		for i := 1; i <= maxRetries; i++ {
			if _, err := bundleService.SaveBundle(context.Background(), request); err != nil {
				l.Warn().
					Str("bundle_name", request.Bundle.Name).
					Str("attempt", fmt.Sprintf("%v/%v", strconv.Itoa(i), strconv.Itoa(maxRetries))).
					Msgf("error creating bundle")
				continue
			} else {
				l.Info().
					Str("bundle_name", request.Bundle.Name).
					Str("after", fmt.Sprintf("%v retries", strconv.Itoa(i))).
					Str("bundleName", request.Bundle.Name).
					Str("bundleId", request.Bundle.Id).
					Msg("default settings bundle registered")
				goto OUT
			}
		}
		l.Err(err).Str("setting_name", request.Bundle.Name).Msg("bundle could not be registered. max number of retries reached")
	} else {
		l.Info().
			Str("bundleName", response.Bundle.Name).
			Str("bundleId", response.Bundle.Id).
			Msg("default settings bundle registered")
	}

OUT:
	permissionRequests := []*settings.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id: "d5f42c4b-e1b6-4b59-8eca-fc4b9e9f2320",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingIDGreeterPhrase,
				},
				Name: "phrase-admin-read",
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READWRITE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}

	for i := range permissionRequests {
		l.Debug().Str("setting_name", permissionRequests[i].Setting.Name).Str("bundle_id", permissionRequests[i].BundleId).Msg("registering setting to bundle")
		if res, err := bundleService.AddSettingToBundle(context.Background(), permissionRequests[i]); err != nil {
			go retryPermissionRequests(context.Background(), bundleService, permissionRequests[i], maxRetries, l)
		} else {
			l.Info().Str("setting_name", res.Setting.Name).Msg("permission registered")
		}
	}
}

// proposal: the retry logic should live in the settings service.
func retryPermissionRequests(ctx context.Context, bs settings.BundleService, setting *settings.AddSettingToBundleRequest, maxRetries int, l *log.Logger) {
	for i := 1; i < maxRetries; i++ {
		if _, err := bs.AddSettingToBundle(ctx, setting); err != nil {
			l.Warn().Str("setting_name", setting.Setting.Name).Str("attempt", fmt.Sprintf("%v/%v", strconv.Itoa(i), strconv.Itoa(maxRetries))).Msgf("error on add setting to bundle")
			continue
		}
		l.Info().Str("setting_name", setting.Setting.Name).Str("after", fmt.Sprintf("%v retries", strconv.Itoa(i))).Msg("permission registered")
		return
	}

	l.Error().Str("setting_name", setting.Setting.Name).Msg("setting could not be registered. max number of retries reached")
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
		value, ok := response.Value.Value.Value.(*settings.Value_StringValue)
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
