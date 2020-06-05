package svc

import (
	"context"
	"errors"
	"fmt"
	mclient "github.com/micro/go-micro/v2/client"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

var (
	// ErrMissingName defines the error if name is missing.
	ErrMissingName = errors.New("missing a name")
)

// NewService returns a service implementation for HelloHandler.
func NewService() v0proto.HelloHandler {
	return Hello{}
}

// Hello defines implements the business logic for HelloHandler.
type Hello struct {
	// Add database handlers here.
}

// Greet implements the HelloHandler interface.
func (s Hello) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	if req.Name == "" {
		return ErrMissingName
	}

	rsp.Message = fmt.Sprintf(
		"Hello %s",
		req.Name,
	)

	return nil
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	request := &settings.SaveSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Identifier: &settings.Identifier{
				Extension: "ocis-hello",
				BundleKey: "greeting",
			},
			DisplayName: "Greeting",
			Settings: []*settings.Setting{
				{
					SettingKey:  "phrase",
					DisplayName: "Phrase",
					Description: "Phrase for replies on the greet request",
					Value: &settings.Setting_StringValue{
						StringValue: &settings.StringSetting{
							Required:  true,
							Default:   "Hello",
							MaxLength: 15,
						},
					},
				},
			},
		},
	}

	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	bundleService := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)
	response, err := bundleService.SaveSettingsBundle(context.Background(), request)
	if err != nil {
		l.Err(err).
			Msg("Error registering settings bundle")
	} else {
		l.Info().
			Str("bundle key", response.SettingsBundle.Identifier.BundleKey).
			Msg("Successfully registered settings bundle")
	}
}
