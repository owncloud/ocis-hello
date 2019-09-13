---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 30
---

As this project is built with Go and NodeJS, so you need to install that first. The installation of Go and NodeJS is out of the scope of this document, please follow the official documentation for [Go](https://golang.org/doc/install), [NodeJS](https://nodejs.org/en/download/package-manager/) and [Yarn](https://yarnpkg.com/lang/en/docs/install/). After the installation of the required tools you need to get the sources:

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis-hello.git
cd ocis-hello
{{< / highlight >}}

All required tool besides Go itself and make are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile` and respectively our `package.json`.

{{< highlight txt >}}
# install dependecies
yarn install

# build frontend
yarn build

# generate assets
make generate

# build binary
make build
{{< / highlight >}}

Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis-hello -h` to see all available options.
