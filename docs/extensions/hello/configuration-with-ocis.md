---
title: "Running"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: configuration-with-ocis.md
---

#### Running ocis
In order to run this extension we will need to run oCIS first. For that clone and build the oCIS single binary from the github repo `https://github.com/owncloud/ocis`.
After that we will need to create a config file for phoenix so that we can load the hello app in the frontend. Create a file `web-config.json` with the following contents.
```json
{
  "server": "https://localhost:9200",
  "theme": "owncloud",
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
Here we can add the url for the js file from where the hello app will be loaded.

After that we will need a configuration file for ocis where we can specify the path for the hello app in the backend. For this you can use the existing `proxy-example.json` file from the [ocis-proxy](https://github.com/owncloud/ocis/blob/master/proxy/config/proxy-example.json) repo. Just add an extra endpoint at the end for the hello app.
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

In addition to all these we will also need to set the config files we just modified. For that set these variables with the path to the config files.
```
export WEB_UI_CONFIG=<path to web-config.json>
export PROXY_CONFIG_FILE=<path to ocis proxy config file>
```
And finally start the ocis server
```
ocis server
```

After this we will need to start the ocis-hello service.
For that just build ocis-hello binary.
```
cd ocis-hello 
make
```
And Run the service
```
bin/ocis-hello server
```
