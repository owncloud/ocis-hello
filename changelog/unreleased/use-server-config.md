Change: Use server config variable from ocis-web

We are not providing an api url anymore but use the server url from ocis-web config instead. This still - as before - requires that ocis-proxy is in place for routing API requests to ocis-hello.

https://github.com/owncloud/ocis-hello/pull/81
