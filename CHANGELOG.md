# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-hello unreleased.

[unreleased]: https://github.com/owncloud/ocis-hello/compare/v0.1.0...master

## Summary

* Bugfix - Fix build error: [#72](https://github.com/owncloud/ocis-hello/pull/72)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#84](https://github.com/owncloud/ocis-hello/pull/84)
* Bugfix - Fix api path: [#99](https://github.com/owncloud/ocis-hello/pull/99)
* Change - Update micro: [#49](https://github.com/owncloud/ocis-hello/pull/49)
* Change - Use server config variable from ocis-web: [#81](https://github.com/owncloud/ocis-hello/pull/81)
* Change - Adapt to new ocis-settings data model: [#85](https://github.com/owncloud/ocis-hello/pull/85)
* Change - Update ocis-pkg and ocis settings: [#593](https://github.com/owncloud/ocis/pull/593)
* Enhancement - Track tool dependencies: [#51](https://github.com/owncloud/ocis-hello/pull/51)
* Enhancement - Streamline project structure: [#79](https://github.com/owncloud/ocis-hello/pull/79)
* Enhancement - Update JS dependencies: [#120](https://github.com/owncloud/ocis-hello/pull/120)

## Details

* Bugfix - Fix build error: [#72](https://github.com/owncloud/ocis-hello/pull/72)

   We had the issue of the flagset being called twice (on both http and grpc server), which cause a
   clash with already defined flags. We are now only calling the flagset once to get back to a
   working build.

   https://github.com/owncloud/ocis-hello/pull/72


* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#84](https://github.com/owncloud/ocis-hello/pull/84)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-hello/pull/84


* Bugfix - Fix api path: [#99](https://github.com/owncloud/ocis-hello/pull/99)

   The server path coming from ownCloud Web now has an enforced trailing slash. Concatenating the
   api path to the server path resulted in a path containing a double slash.

   https://github.com/owncloud/ocis-hello/pull/99


* Change - Update micro: [#49](https://github.com/owncloud/ocis-hello/pull/49)

   Updated the micro dependencies.

   https://github.com/owncloud/ocis-hello/pull/49


* Change - Use server config variable from ocis-web: [#81](https://github.com/owncloud/ocis-hello/pull/81)

   We are not providing an api url anymore but use the server url from ocis-web config instead. This
   still - as before - requires that ocis-proxy is in place for routing API requests to ocis-hello.

   https://github.com/owncloud/ocis-hello/pull/81


* Change - Adapt to new ocis-settings data model: [#85](https://github.com/owncloud/ocis-hello/pull/85)

   Ocis-settings introduced UUIDs and less verbose endpoint and message type names. This PR
   adjusts ocis-hello accordingly.

   https://github.com/owncloud/ocis-hello/pull/85
   https://github.com/owncloud/ocis-settings/pull/46


* Change - Update ocis-pkg and ocis settings: [#593](https://github.com/owncloud/ocis/pull/593)

   Ocis-pkg and ocis settings have been moved to the ocis mono-repo.

   https://github.com/owncloud/ocis/pull/593
   https://github.com/owncloud/ocis-hello/pull/99


* Enhancement - Track tool dependencies: [#51](https://github.com/owncloud/ocis-hello/pull/51)

   Added tracking for tool dependencies to be able to run go mod tidy without losing them. More
   information:
   https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

   https://github.com/owncloud/ocis-hello/pull/51


* Enhancement - Streamline project structure: [#79](https://github.com/owncloud/ocis-hello/pull/79)

   - We have aligned the project structure of ocis-hello with other repositories and improved
   error logging. - When running this service through `make watch` it now regenerates embedded
   assets properly as soon as the web bundle is changed / saved. - In the package.json file we're now
   declaring owncloud-design-system as peer dependency, since we're actively using it. It
   comes from ocis-web, so we don't need to bundle it.

   https://github.com/owncloud/ocis-hello/pull/79
   https://github.com/owncloud/ocis-hello/pull/80


* Enhancement - Update JS dependencies: [#120](https://github.com/owncloud/ocis-hello/pull/120)

   We've bumped the JS dependencies (including the ownCloud design system) and removed an unused
   `ldap` package.

   https://github.com/owncloud/ocis-hello/pull/120

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

