package svc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	mclient "github.com/micro/go-micro/v2/client"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis-settings/pkg/service/v0"
)

var (
	// ErrMissingName defines the error if name is missing.
	ErrMissingName = errors.New("missing a name")

	bundleIdGreeting       = "21fb587b-7b69-4aa6-b0a7-93c74af1918f"
	settingIdGreeterPhrase = "b3584ea8-caec-4951-a2c1-92cbc70071b7"
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

	phrase := getGreetingPhrase(ctx)
	rsp.Message = fmt.Sprintf(phrase, req.Name)

	return nil
}

func getGreetingPhrase(ctx context.Context) string {
	ownAccountUUID := ctx.Value(middleware.UUIDKey)
	if ownAccountUUID != nil {
		// request to the settings service requires to have the account uuid of the authenticated user available in the context
		request := &settings.GetValueRequest{
			Identifier: &settings.Identifier{
				Extension:   "ocis-hello",
				BundleKey:   "greeting",
				SettingKey:  "phrase",
				AccountUuid: ownAccountUUID.(string),
			},
		}

		// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
		// https://github.com/owncloud/ocis-hello/issues/74
		valueService := settings.NewValueService("com.owncloud.api.settings", mclient.DefaultClient)
		response, err := valueService.GetValue(ctx, request)
		if err == nil {
			value := response.Value.Value.Value.(*settings.Value_StringValue)
			trimmedPhrase := strings.Trim(
				value.StringValue,
				" \t",
			)
			if trimmedPhrase != "" {
				return trimmedPhrase + " %s"
			}
		}
	}
	return "Hello %s"
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	request := &settings.SaveBundleRequest{
		Bundle: &settings.Bundle{
			Id:          bundleIdGreeting,
			Name:        "greeting",
			DisplayName: "Greeting",
			Extension:   "ocis-hello",
			Type:        settings.Bundle_TYPE_DEFAULT,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			Settings: []*settings.Setting{
				{
					Id:          settingIdGreeterPhrase,
					Name:        "phrase",
					DisplayName: "Phrase",
					Description: "Phrase for replies on the greet request",
					Resource: &settings.Resource{
						Type: settings.Resource_TYPE_USER,
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

	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	bundleService := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)
	response, err := bundleService.SaveBundle(context.Background(), request)
	if err != nil {
		l.Err(err).
			Msg("Error registering settings bundle")
	} else {
		l.Info().
			Str("bundleName", response.Bundle.Name).
			Str("bundleId", response.Bundle.Id).
			Msg("Successfully registered settings bundle")
	}

	permissionRequests := []*settings.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id: "45f52511-f9cc-4226-92e3-779e07179714",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingIdGreeterPhrase,
				},
				Name: "phrase-user-read",
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READ,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id: "d5f42c4b-e1b6-4b59-8eca-fc4b9e9f2320",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingIdGreeterPhrase,
				},
				Name: "phrase-admin-read",
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READ,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id: "8732811a-147c-4b28-89f5-112573c40682",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingIdGreeterPhrase,
				},
				Name: "phrase-admin-create",
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_CREATE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id: "9bd896e2-127e-4946-871d-3d1c5f2a52f2",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingIdGreeterPhrase,
				},
				Name: "phrase-admin-update",
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_UPDATE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}


}
