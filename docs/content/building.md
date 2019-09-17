---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 30
---

## Project Dependencies

The project has the following dependencies:

- [GNU Make](gnu-make)
- [Go](golang),
- [NodeJS](nodejs),
- [Yarn](yarn),

## Get The Source

After installing all of the project’s dependencies, you need to get the source, and switch to the cloned source directory. The example below shows how.

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis-hello.git
cd ocis-hello
{{< / highlight >}}

**Note:** All required Go packages are installed inside [the GOPATH directory][gopath].

## Frontend

The commands below install the required build dependencies and build the frontend bundle. This frontend bundle will be embedded into the generated backend binary in the next step.

{{< highlight txt >}}
yarn install
yarn build
{{< / highlight >}}

## Backend

The commands below embed the generated frontend bundle into the generated backend binary.

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

If the commands above complete without error, the backend binary will be available in the `bin/` folder. You can test it by running `./bin/ocis-hello -h`.

[gnu-make]: https://www.gnu.org/software/make/
[golang]: https://golang.org/doc/install
[gopath]: https://github.com/golang/go/wiki/GOPATH
[nodejs]: https://nodejs.org/en/download/package-manager/
[yarn]: https://yarnpkg.com/lang/en/docs/install/
