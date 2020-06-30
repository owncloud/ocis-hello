# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-hello unreleased.

[unreleased]: https://github.com/owncloud/ocis-hello/compare/v0.1.0...master

## Summary

* Bugfix - Fix build error: [#72](https://github.com/owncloud/ocis-hello/pull/72)
* Change - Update micro: [#49](https://github.com/owncloud/ocis-hello/pull/49)
* Enhancement - Streamline project structure: [#79](https://github.com/owncloud/ocis-hello/pull/79)
* Enhancement - Track tool dependencies: [#51](https://github.com/owncloud/ocis-hello/pull/51)

## Details

* Bugfix - Fix build error: [#72](https://github.com/owncloud/ocis-hello/pull/72)

   We had the issue of the flagset being called twice (on both http and grpc server), which cause a
   clash with already defined flags. We are now only calling the flagset once to get back to a
   working build.

   https://github.com/owncloud/ocis-hello/pull/72


* Change - Update micro: [#49](https://github.com/owncloud/ocis-hello/pull/49)

   Updated the micro dependencies.

   https://github.com/owncloud/ocis-hello/pull/49


* Enhancement - Streamline project structure: [#79](https://github.com/owncloud/ocis-hello/pull/79)

   We have aligned the project structure of ocis-hello with other repositories and improved
   error logging.

   https://github.com/owncloud/ocis-hello/pull/79


* Enhancement - Track tool dependencies: [#51](https://github.com/owncloud/ocis-hello/pull/51)

   Added tracking for tool dependencies to be able to run go mod tidy without losing them. More
   information:
   https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

   https://github.com/owncloud/ocis-hello/pull/51

# Changelog for [0.1.0] (2020-01-24)

The following sections list the changes in ocis-hello 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-hello/compare/c43f3a33cb0b57d7e25ebc88c138d22e95f88cfe...v0.1.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-hello/issues/1)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-hello/issues/1)

   Just prepared an initial basic version to serve a hello world API that also provides a Phoenix
   extension to demonstrate the plugin architecture in combination with Phoenix and ownCloud
   Infinite Scale.

   https://github.com/owncloud/ocis-hello/issues/1

