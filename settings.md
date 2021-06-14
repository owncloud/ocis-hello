---
title: "Settings"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: settings.md
---

oCIS provides an settings extension that other extensions can use to make settings accessible to users.

In case of the Hello extension administrator users are able to change the greeter message.

Settings are stored and presented to the user by the oCIS settings extension. It also exposes endpoints for registering and manipulating so called
[settings bundles](https://owncloud.dev/extensions/settings/bundles/).

The settings defined through settings bundles can be changed by authenticated users in ownCloud Web, if they have enough permissions to edit them. As a result, your service only has to register a settings bundle and permissions for it and oCIS settings takes care of everything else. Your service can simply use the settings values that were set by users.

In this chapter we want to show you how to register a settings bundle, the settings permissions and how to use the respective values that were set by users. We do this by customizing the greeter phrase from our greeter service in the Hello extension.

You can find the source code, especially how it's integrated into the service, in the following files:
- `pkg/service/v0/service.go` for the requests
- `pkg/command/server.go` for the integration of `registerSettingsBundles` into the service start

## Register a settings bundle and set the permissions
In order to register a settings bundle, you need to create a request message and then send it
to the `BundleService` of `oCIS settings` through a gRPC call. The same applies for setting permissions on the setting bundles.

### Create a bundle request

```go
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

```
The request holds only one field, which is a `SettingsBundle`. It consists of an `Identifier`, a `DisplayName`
and a list of `Settings`.
- The `Extension` and the `ID` of the bundle are required and have to be
alphanumeric (`-` and `_` are allowed as well). The `ID` has to stay stable - if you change it, existing settings will not be migrated to the new identifier.
- The `DisplayName` is required and may contain any UTF8 character. It will be shown in the settings user frontend in a generated form, so please try to be descriptive. You can change the `DisplayName` at any time.
- `Settings` is the list of settings you want to make available with this settings bundle. In this example, there
is only one setting defined - a string setting for the phrase our greeter uses in the response. You can explore
more types of settings in the `settings` package. All of them come with their own characteristics and validations. For the `phrase` setting we decided to set it to `Required`, so that it can't be empty, and to set a `MaxLength` of 15 characters, so that the phrase is not too long. The `ID` of the setting is again particularly important, as this is used for referencing the setting in other requests. It has to fulfill the same rules as the other `ID` attributes. Please also take the time to set a `Description`, in order to provide accessibility in the generated forms as good as possible.

### Send bundle request to oCIS settings
This request message can be sent to the `BundleService` of `oCIS settings` like this:
```go
settings := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)
response, err := bundleService.SaveBundle(context.Background(), request)
```

We run this request on every start of the Hello extension so that the settings service always has the most recent version of the settings bundle.

### Create a permission settings request
In order to grant admins access to the setting we need to create a `AddSettingToBundleRequest`.

```go
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
```

### Send permission settings request
The `AddSettingToBundleRequest` then needs to be send to the settings service:

```go
response, err := bundleService.AddSettingToBundle(context.Background(), permissionRequests[i])
```

## Use settings value

We registered the greeter phrase setting for a reason: We want to allow authenticated administrator users to customize how users are greeted by the Hello extension. In order to do this, we need to ask oCIS settings on every request, what the greeter phrase of the authenticated user is.

### Account UUID
The settings request has one important prerequisite: As our service is stateless, we need to know the account UUID of the authenticated user the incoming POST request to our greeter service is coming from.
As that POST request is coming through oCIS proxy, there is an HTTP header `x-access-token` that holds a JWT with the account UUID in it. We just have to dismantle the JWT to get the UUID. There is a middleware for that in `ocis-pkg`. You can look up the server configuration for that middleware in `pkg/server/http/server.go`.
In essence, it dismantles the `x-access-token`, extracts the account UUID and makes it available in the context.

It can be subsequently retrieved from the context like this:
```go
accountID, ok := metadata.Get(ctx, middleware.AccountID)
```

### Create request
With the account UUID we can build an `Identifier` for the request to oCIS settings as follows:
```go
rq := settings.GetValueByUniqueIdentifiersRequest{
  AccountUuid: accountID,
  SettingId:   settingIDGreeterPhrase,
}

response, err := s.vsClient.GetValueByUniqueIdentifiers(context.Background(), &rq)
```
In order to get the setting we need to know which user (`AccountUuid`) is requesting the settings value and which setting he is requesting (`SettingId`).

### Send request
This request message can be sent to the settings extension like this:

```go
valueService := settings.NewValueService("com.owncloud.api.settings", mclient.DefaultClient)
response, err := valueService.GetSettingsValue(ctx, request)
```

The request gives us the default or customized greeting phrase, depending on wether it has already been changed by an administrator.

## Conclusion
You have learned how to register *settings bundles*, how to get the account UUID of the authenticated user and how to query the settings service for *settings values*.
