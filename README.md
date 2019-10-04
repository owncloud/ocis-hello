# ownCloud Infinite Scale: Hello

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-hello/status.svg)](https://cloud.drone.io/owncloud/ocis-hello)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6f1eaaa399294d959ef7b3b10deed41d)](https://www.codacy.com/manual/owncloud/ocis-hello?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-hello&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-hello?status.svg)](http://godoc.org/github.com/owncloud/ocis-hello)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-hello)](http://goreportcard.com/report/github.com/owncloud/ocis-hello)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fowncloud%2Focis-hello.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fowncloud%2Focis-hello?ref=badge_shield)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-hello.svg)](http://microbadger.com/images/owncloud/ocis-hello "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/hello/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.12. For the frontend it's also required to have [NodeJS](https://nodejs.org/en/download/package-manager/) and [Yarn](https://yarnpkg.com/lang/en/docs/install/) installed.

```console
git clone https://github.com/owncloud/ocis-hello.git
cd ocis-hello

yarn install
yarn build

make generate build

./bin/ocis-hello -h
```

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fowncloud%2Focis-hello.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fowncloud%2Focis-hello?ref=badge_large)

## Copyright

```console
Copyright (c) 2019 ownCloud GmbH <https://owncloud.com>
```
