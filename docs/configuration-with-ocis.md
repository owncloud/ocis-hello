---
title: "Running"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: configuration-with-ocis.md
---

### Configuring ocis-hello with ocis
We will need various services to run ocis

#### Running a ldap server in docker container
We will use the ldap server as users provider for ocis.
```
docker run --hostname ldap.my-company.com \
    -e LDAP_TLS_VERIFY_CLIENT=never \
    -e LDAP_DOMAIN=owncloud.com \
    -e LDAP_ORGANISATION=ownCloud \
    -e LDAP_ADMIN_PASSWORD=admin \
    --name docker-slapd \
    -p 127.0.0.1:389:389 \
    -p 636:636 -d osixia/openldap
```
#### Running a redis server in a docker container
Redis will be used by ocis for various caching purposes.
```
docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest

```
#### Running ocis
In order to run this extension we will need to run ocis first. For that clone and build the ocis single binary from the github repo `https://github.com/owncloud/ocis`.
After that we will need to create a config file for phoenix so that we can load the hello app in the frontend. Create a file `phoenix-config.json` with the following contents.
```json
{
    "server": "https://localhost:9200",
    "theme": "owncloud",
    "version": "0.1.0",
    "openIdConnect": {
        "metadata_url": "https://localhost:9200/.well-known/openid-configuration",
        "authority": "https://localhost:9200",
        "client_id": "phoenix",
        "response_type": "code",
        "scope": "openid profile email"
    }, 
    "apps": [
        "files",
        "draw-io",
        "pdf-viewer",
        "markdown-editor",
        "media-viewer"
    ], 
    "external_apps": [
        {
            "id": "hello",
            "path": "http://localhost:9105/hello.js",
            "config": {
                "url": "http://localhost:9105"
            }
        }
    ]   
}
```
Here we can add the url for the js file from where the hello app will be loaded.

After that we will need a configuration file for ocis where we can specify the path for the hello app in the backend. For this you can use the existing `proxy-example.json` file from the [ocis-proxy](https://github.com/owncloud/ocis-proxy/blob/master/config/proxy-example.json) repo. Just add an extra endpoint at the end for the hello app.
```json
        {
          "endpoint": "/api/v0/greet",
          "backend": "http://localhost:9105"
        }
```

With all this in place we can finally start ocis. But first we will need to set some configuration variables.
```
export REVA_USERS_DRIVER=ldap
export REVA_LDAP_HOSTNAME=localhost
export REVA_LDAP_PORT=636
export REVA_LDAP_BASE_DN='dc=owncloud,dc=com'
export REVA_LDAP_USERFILTER='(&(objectclass=posixAccount)(cn=%s))'
export REVA_LDAP_GROUPFILTER='(&(objectclass=posixGroup)(cn=%s))'
export REVA_LDAP_BIND_DN='cn=admin,dc=owncloud,dc=com'
export REVA_LDAP_BIND_PASSWORD=admin
export REVA_LDAP_SCHEMA_UID=uid
export REVA_LDAP_SCHEMA_MAIL=mail
export REVA_LDAP_SCHEMA_DISPLAYNAME=displayName
export REVA_LDAP_SCHEMA_CN=cn
export LDAP_URI=ldap://localhost
export LDAP_BINDDN='cn=admin,dc=owncloud,dc=com'
export LDAP_BINDPW=admin
export LDAP_BASEDN='dc=owncloud,dc=com'
```

In addition to all these we will also need to set the config files we just modified. For that set these variables with the path to the config files.
```
export PHOENIX_WEB_CONFIG=<path to phoenix config file>
export OCIS_CONFIG_FILE=<path to ocis proxy config file>
```
And finally start the ocis server
```
bin/ocis server
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
