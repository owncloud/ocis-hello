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
  "version": "0.1.0",
  "openIdConnect": {
    "metadata_url": "https://localhost:9200/.well-known/openid-configuration",
    "authority": "https://localhost:9200",
    "client_id": "web",
    "response_type": "code",
    "scope": "openid profile email"
  },
  "apps": ["files", "media-viewer"],
  "external_apps": [
    {
      "id": "settings",
      "path": "/settings.js"
    },
    {
      "id": "accounts",
      "path": "/accounts.js"
    },
    {
      "id": "hello",
      "path": "/hello.js"
    }
  ],
  "options": {
    "hideSearchBar": true
  }
}

```

Please note the the regististration of our Hello extension in the `external_apps` section. It will trigger ownCloud Web to load `hello.js`, the frontend bundle generated in the [frontend build step]({{< ref "./building#frontend">}}).

The frontend bundle will be requested from the oCIS proxy and requests to our Hello extension's API will also be passed to the oCIS proxy first. Therefore the oCIS proxy needs to be configured to forward these requests to our Hello extension.
We do this by using the existing `proxy-example.json` file from the [oCIS proxy](https://github.com/owncloud/ocis/blob/master/proxy/config/proxy-example.json). Just add an extra endpoint at the end for the Hello extension.

```json
{
  "endpoint": "/api/v0/greet",
  "backend": "http://localhost:9105"
},
{
  "endpoint": "/hello.js",
  "backend": "http://localhost:9105"
}
```

In addition to all these we will also need to activate the config files we just created. Therefore set these two variables with the path to the respective config file.
```
export WEB_UI_CONFIG=<path to web-config.json>
export PROXY_CONFIG_FILE=<path to ocis proxy config file>
```
And finally start the oCIS server:
```
ocis server
```

{{< hint warning >}}
oCIS has currently a bug, that oCIS proxy will not pick up the proxy configuration file if it is started in the supervised mode by `ocis server`. Therefore you will need to apply following workaround:

Run `ocis server` with the environment variables mentioned above. They open a new CLI and run `ocis kill proxy`. Set the same environment variables as above and run `ocis proxy server`. This starts the proxy in a non supervised mode and ensures that it picks up your custom routes in the proxy configuration file.
{{< /hint >}}

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
