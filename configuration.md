---
title: "Configuration"
date: "2021-08-30T17:23:42+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/templates
geekdocFilePath: CONFIGURATION.tmpl
---

{{< toc >}}

## Configuration

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-hello reads `hello.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis-hello/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command line flags

If you prefer to configure the service with command line flags you can see the available variables below. Command line flags are only working when calling the sub command directly.

## Root Command

Hello, an example oCIS extension

Usage: `hello [global options] command [command options] [arguments...]`


-config-file |  $HELLO_CONFIG_FILE
: Path to config file.


-log-level |  $HELLO_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level. Default: `info`.


-log-pretty |  $HELLO_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $HELLO_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





















## Sub Commands

### hello health

Check health status

Usage: `hello health [command options] [arguments...]`






-debug-addr |  $HELLO_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9109`.




















### hello server

start hello service

Usage: `hello server [command options] [arguments...]`







-log-file |  $HELLO_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-tracing-enabled |  $HELLO_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $HELLO_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.


-tracing-endpoint |  $HELLO_TRACING_ENDPOINT
: Endpoint for the agent.


-tracing-collector |  $HELLO_TRACING_COLLECTOR
: Endpoint for the collector.


-tracing-service |  $HELLO_TRACING_SERVICE
: Service name for tracing. Default: `hello`.


-debug-addr |  $HELLO_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9109`.


-debug-token |  $HELLO_DEBUG_TOKEN
: Token to grant metrics access.


-debug-pprof |  $HELLO_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $HELLO_DEBUG_ZPAGES
: Enable zpages debugging.


-http-namespace |  $HELLO_HTTP_NAMESPACE
: Set the base namespace for the http namespace. Default: `com.owncloud.web`.


-http-addr |  $HELLO_HTTP_ADDR
: Address to bind http server. Default: `0.0.0.0:9105`.


-http-root |  $HELLO_HTTP_ROOT
: Root path of http server. Default: `/`.


-http-cache-ttl |  $HELLO_CACHE_TTL
: Set the static assets caching duration in seconds. Default: `604800`.


-grpc-namespace |  $HELLO_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.


-name |  $HELLO_NAME
: service name. Default: `"hello"`.


-grpc-addr |  $HELLO_GRPC_ADDR
: Address to bind grpc server. Default: `0.0.0.0:9106`.


-asset-path |  $HELLO_ASSET_PATH
: Path to custom assets.


-jwt-secret |  $HELLO_JWT_SECRET
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

