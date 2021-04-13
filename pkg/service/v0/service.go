package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/owncloud/ocis-hello/pkg/config"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis/settings/pkg/service/v0"
)

var (
	// ErrMissingName defines the error if name is missing.
	ErrMissingName = errors.New("missing a name")

	bundleIDGreeting       = "21fb587b-7b69-4aa6-b0a7-93c74af1918f"
	settingIDGreeterPhrase = "b3584ea8-caec-4951-a2c1-92cbc70071b7"

	// maxRetries indicates how many times to try a request for network reasons.
	maxRetries = 5
)

// New returns a new instance of Service
func New(opts ...Option) (s *Service, err error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	bundleService := options.BundleService
	if bundleService == nil {
		bundleService = settings.NewBundleService("com.owncloud.api.settings", grpc.DefaultClient)
	}

	s = &Service{
		id:            cfg.GRPC.Namespace + "." + cfg.Server.Name,
		log:           logger,
		Config:        cfg,
		BundleService: bundleService,
	}

	registerSettingsBundles(s, &logger)

	return s, nil
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	id            string
	log           log.Logger
	Config        *config.Config
	BundleService settings.BundleService
}

// Greet implements the HelloHandler interface.
func (s Service) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	if req.Name == "" {
		return ErrMissingName
	}

	phrase := getGreetingPhrase(ctx)
	rsp.Message = fmt.Sprintf(phrase, req.Name)

	return nil
}

func getGreetingPhrase(ctx context.Context) string {
	ownAccountUUID := ctx.Value(middleware.UUIDKey)
	if ownAccountUUID != nil {
		// request to the settings service requires to have the account uuid of the authenticated user available in the context
		rq := settings.GetValueByUniqueIdentifiersRequest{
			AccountUuid: ownAccountUUID.(string),
			SettingId:   settingIDGreeterPhrase,
		}

		valueService := settings.NewValueService("com.owncloud.api.settings", grpc.DefaultClient)
		response, err := valueService.GetValueByUniqueIdentifiers(ctx, &rq)
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
	}
	return "Hello %s"
}

// registerSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func registerSettingsBundles(s *Service, l *log.Logger) {
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

	response, err := s.BundleService.SaveBundle(context.Background(), request)
	if err != nil {
		l.Warn().Msg("error registering settings bundle at first try. retrying")
		for i := 1; i <= maxRetries; i++ {
			if _, err := s.BundleService.SaveBundle(context.Background(), request); err != nil {
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
		if res, err := s.BundleService.AddSettingToBundle(context.Background(), permissionRequests[i]); err != nil {
			go retryPermissionRequests(context.Background(), s.BundleService, permissionRequests[i], maxRetries, l)
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
