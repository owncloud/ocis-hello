---
title: "Building"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: building.md
---

{{< toc >}}

As this project is built with Go and NodeJS, so you need to install that first. The installation of Go and NodeJS is out of the scope of this document, please follow the official documentation for [Go](https://golang.org/doc/install), [NodeJS](https://nodejs.org/en/download/package-manager/) and [Yarn](https://yarnpkg.com/lang/en/docs/install/), to build this project you have to install Go >= v1.16. After the installation of the required tools you need to get the sources:

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis-hello.git
cd ocis-hello
{{< / highlight >}}

All required tool besides the ones mentioned above will be automatically installed. All commands to build this project are part of our `Makefile` and respectively our `package.json`.

## Frontend

{{< highlight txt >}}
yarn install
yarn build
{{< / highlight >}}

The above commands will install the required build dependencies and build the whole frontend bundle. This bundle will we embedded into the binary later on.

## Backend

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

The above commands will embed the frontend bundle into the binary. Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/hello -h` to see all available options.

## Documentation

Just run `make -C docs docs-serve` from within the root level of the extensions git repository. This will make documentation available on [localhost:1313](http://localhost:1313) and also do a hot reload if you change something in the (non autogenerated) documentation files.