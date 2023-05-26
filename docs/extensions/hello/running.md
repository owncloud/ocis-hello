---
title: "Running"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: running.md
---

In order to use the Hello extension we need to configure and start oCIS first. After that we can run the Hello extension itself.

## Configure and start oCIS
You can either [start oCIS from a binary](https://owncloud.dev/ocis/getting-started/#binaries) or [build it from source](https://owncloud.dev/ocis/development/build/#build-the-ocis-binary).

No matter which way you choose, we need to create a configuration file for ownCloud Web, so that the Hello extension will be loaded in the frontend. Therefore create a file `web-config.json` with the following contents:
```json
{
  "server": "https://localhost:9200",
  "theme": "https://localhost:9200/themes/owncloud/theme.json",
  "openIdConnect": {
    "metadata_url": "https://localhost:9200/.well-known/openid-configuration",
    "authority": "https://localhost:9200",
    "client_id": "web",
    "response_type": "code",
    "scope": "openid profile email"
  },
  "apps": [
    "files",
    "preview",
    "pdf-viewer",
    "search",
    "text-editor",
    "draw-io",
    "external",
    "admin-settings"
  ],
  "external_apps": [
    {
      "id": "hello",
      "path": "/hello.js"
    }
  ],
  "options": {
    "previewFileMimeTypes": [
      "image/gif",
      "image/png",
      "image/jpeg",
      "text/plain",
      "image/tiff",
      "image/bmp",
      "image/x-ms-bmp"
    ],
    "upload": {
      "xhr": {
        "timeout": 60000
      }
    },
    "contextHelpersReadMore": true
  }
}
```

Please note the the regististration of our Hello extension in the `external_apps` section. It will trigger ownCloud Web to load `hello.js`, the frontend bundle generated in the [frontend build step]({{< ref "./building#frontend">}}).

To activate the config file we just created we need to set this variable with the path to it.
```
export WEB_UI_CONFIG=<path to web-config.json>
```

The frontend bundle will be requested from the oCIS proxy and requests to our Hello extension's API will also be passed to the oCIS proxy first. Therefore the oCIS proxy needs to be configured to forward these requests to our Hello extension.
In the ocis config folder create a file called `proxy.yaml` with this content:

```yaml
policies:
- name: ocis
  routes:
    - endpoint: /
      service: com.owncloud.web.web
      unprotected: true
    - endpoint: /.well-known/webfinger
      service: com.owncloud.web.webfinger
      unprotected: true
    - endpoint: /.well-known/openid-configuration
      service: com.owncloud.web.idp
      unprotected: true
    - endpoint: /branding/logo
      service: com.owncloud.web.web
    - endpoint: /konnect/
      service: com.owncloud.web.idp
      unprotected: true
    - endpoint: /signin/
      service: com.owncloud.web.idp
      unprotected: true
    - endpoint: /archiver
      service: com.owncloud.web.frontend
    - endpoint: /ocs/v2.php/apps/notifications/api/v1/notifications
      service: com.owncloud.userlog.userlog
    - type: regex
      endpoint: /ocs/v[12].php/cloud/user/signing-key
      service: com.owncloud.web.ocs
    - type: regex
      endpoint: /ocs/v[12].php/config
      service: com.owncloud.web.frontend
      unprotected: true
    - endpoint: /ocs/
      service: com.owncloud.web.frontend
    - type: query
      endpoint: /remote.php/?preview=1
      service: com.owncloud.web.webdav
    - method: REPORT
      endpoint: /remote.php/dav/
      service: com.owncloud.web.webdav
    - method: REPORT
      endpoint: /remote.php/webdav
      service: com.owncloud.web.webdav
    - method: REPORT
      endpoint: /dav/spaces
      service: com.owncloud.web.webdav
    - type: query
      endpoint: /dav/?preview=1
      service: com.owncloud.web.webdav
    - type: query
      endpoint: /webdav/?preview=1
      service: com.owncloud.web.webdav
    - endpoint: /remote.php/
      service: com.owncloud.web.ocdav
    - endpoint: /dav/
      service: com.owncloud.web.ocdav
    - endpoint: /webdav/
      service: com.owncloud.web.ocdav
    - endpoint: /status
      service: com.owncloud.web.ocdav
      unprotected: true
    - endpoint: /status.php
      service: com.owncloud.web.ocdav
      unprotected: true
    - endpoint: /index.php/
      service: com.owncloud.web.ocdav
    - endpoint: /apps/
      service: com.owncloud.web.ocdav
    - endpoint: /data
      service: com.owncloud.web.frontend
      unprotected: true
    - endpoint: /app/list
      service: com.owncloud.web.frontend
      unprotected: true
    - endpoint: /app/
      service: com.owncloud.web.frontend
    - endpoint: /graph/v1.0/invitations
      service: com.owncloud.graph.invitations
    - endpoint: /graph/
      service: com.owncloud.graph.graph
    - endpoint: /api/v0/settings
      service: com.owncloud.web.settings
    - endpoint: /api/v0/greet
      backend: "http://localhost:9105"
    - endpoint: /hello.js
      backend: "http://localhost:9105"
      unprotected: true
```
{{< hint warning >}}
These routes are the default routes of the proxy plus extra routes for the hello extension.
For details see the [proxy documentation](https://owncloud.dev/services/proxy).
{{< /hint >}}

In addition to all these we need to make sure the hello service can be registered to oCIS, for that set this variable.
```
export MICRO_REGISTRY=mdns
```
And finally start the oCIS server:
```
ocis server
```

## Start the extension

After oCIS is running, we can start the Hello extension.

For that just build ocis-hello binary.
```
cd ocis-hello
make build
```
And Run the service
```
./bin/hello server
```
