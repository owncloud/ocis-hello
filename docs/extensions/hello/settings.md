---
title: "Settings"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: settings.md
---

The *ocis-settings* service exposes an endpoint for registering so called
[settings bundles](https://owncloud.github.io/extensions/ocis_settings/bundles/).
This gives control to every service to define settings that are needed for fulfilling it's intended purpose.
There are different types of settings available out of the box - hopefully those already fit your needs.
The settings defined through settings bundles can be changed by authenticated users through an `ocis-web`
extension, which is also provided by the settings service. As a result, your service only has to register
a settings bundle and oCIS takes care of everything else. Your service can simply use the settings values
that were set by users.

With this chapter we want to show you, how to register a settings bundle and how to use the respective
values that were set by users. We do this by customizing the greeter phrase from our greeter service in ocis-hello.

You can find the source code, especially how it's integrated into the service, in the following files:
- `pkg/service/v0/service.go` for the requests,
- `pkg/command/server.go` for the integration of `RegisterSettingsBundles` into the service start.

## Register settings bundle

In order to register a settings bundle, you need to create a request message and then send it
to the `BundleService` of `ocis-settings` through a gRPC call.

### Create request
```go
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
```
The request holds only one field, which is a `SettingsBundle`. It consists of an `Identifier`, a `DisplayName`
and a list of `Settings`.
- The `Extension` and `BundleKey` inside the `Identifier` are required and have to be
alphanumeric (`-` and `_` are allowed as well). The `Identifier` has to stay stable - if you change it, existing
settings will not be migrated to the new identifier.
- The `DisplayName` is required and may contain any UTF8 character. It will be shown in the settings user frontend
in a generated form, so please try to be descriptive. You can change the `DisplayName` at any time.
- `Settings` is the list of settings you want to make available with this settings bundle. In this example, there
is only one setting defined - a string setting for the phrase our greeter uses in the response. You can explore
more types of settings in the `settings` package. All of them come with their own characteristics and validations.
For the `phrase` setting we decided to set it to `Required`, so that it can't be empty, and to set a `MaxLength`
of 15 characters, so that the phrase is not too long. The `SettingKey` is particularly important, as this is
used for referencing the setting in other requests. It has to fulfill the same rules as the other `Identifier`
attributes. Please also take the time to set a `Description`, in order to provide accessibility in the generated
forms as good as possible.

### Send request to ocis-settings
This request message can be sent to the `BundleService` of `ocis-settings` like this:
```go
bundleService := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)
response, err := bundleService.SaveSettingsBundle(context.Background(), request)
```

We run this request on every start of `ocis-hello` so that the settings service always has the most recent
version of the settings bundle.

## Use settings value

We registered the greeter phrase setting for a reason: We want to allow the authenticated user to customize
how they want to be greeted by `ocis-hello`. In order to do this, we need to ask `ocis-settings` on every
request, what the greeter phrase of the authenticated user is.

### Account UUID
The settings request has one important prerequisite: As our service is stateless, we need to know the
account UUID of the authenticated user the incoming POST request to our greeter service is coming from.
As that POST request is coming through `ocis-proxy`, there is an HTTP header `x-access-token` that holds
a JWT with the account UUID in it. We just have to dismantle the JWT to get the UUID. There is a middleware for
that in `ocis-pkg`. You can look up the server configuration for that middleware in `pkg/server/http/server.go`.
In essence, it dismantles the `x-access-token`, extracts the account UUID and makes it available in the context.
It can be subsequently retrieved from the context like this:
```go
ownAccountUUID := ctx.Value(middleware.UUIDKey).(string)
```

### Create request
With the account UUID we can build an `Identifier` for the request to `ocis-settings` as follows:
```go
request := &settings.GetSettingsValueRequest{
    Identifier: &settings.Identifier{
        Extension:   "ocis-hello",
        BundleKey:   "greeting",
        SettingKey:  "phrase",
        AccountUuid: ownAccountUUID,
    },
}
```
The `Identifier` for getting a value from `ocis-settings` needs the three fields for extension, bundle and
setting that we chose in the settings bundle. The fourth field is the UUID of the authenticated user.

### Send request to ocis-settings
This request message can be sent to the `ValueService` of `ocis-settings` like this:
```go
valueService := settings.NewValueService("com.owncloud.api.settings", mclient.DefaultClient)
response, err := valueService.GetSettingsValue(ctx, request)
```

If this request is successful we will have a - possibly customized - greeting phrase. It could also be the
default value, if the user didn't customize their phrase in the settings frontend.

## Conclusion
You have learned how to register *settings bundles*, how to get the account UUID of the authenticated user
and how to query the settings service for *settings values*.
