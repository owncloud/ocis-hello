Bugfix: Fix build error

We had the issue of the flagset being called twice (on both http and grpc server), which cause a clash with already
defined flags. We are now only calling the flagset once to get back to a working build.

https://github.com/owncloud/ocis-hello/pull/72
