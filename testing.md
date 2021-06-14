---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: testing.md
---
This repository provides a general guideline for creating tests for an oCIS extension. The tests can be written in various levels from unit, integration, and end-to-end. It is not essential to write tests on all these levels as it can be redundant in some cases. This repository provides a reference for all levels of tests.

## Unit tests
Unit tests generally live inside *_test.go files in the `/pkg` directory. One such example in this extension is in `/pkg/service/v0/service_test.go`. Similarly the unit test for the protobuf generated code can also be written just like in `/pkg/proto/hello.pb_test.go`.

## Integration tests
There are mainly 2 types of integration tests, namely HTTP tests, and GRPC tests. These tests mostly live in `/pkg/proto` directory where all the protobuf definitions are specified. The examples for the HTTP integration tests are in `/pkg/proto/hello.pb.web_test.go` whereas the GRPC tests are in `/pkg/proto/hello.pb.micro_test.go`.

### End-to-End tests
For extensions with an UI, we can also write end-to-end tests using the Nightwatch test framework. These tests live in `/ui/tests` directory. We can reuse already existing Gherkin steps from the [ownCloud Web](https://github.com/owncloud/web) tests here.

## Running the tests
### Unit and integration tests
The unit and integration tests are run using the simple `go test` command. If you wish to run all the tests with the coverage you can just use make command.
```bash
make test
```
You can also run a specific file with the go test command
```bash
go test <path to package or file>
```
### End-to-End tests
Running end-to-end tests is a bit more complicated than unit and integration tests. First of all we will need a complete oCIS setup with the Hello extension running. Please refer to [foo]({{< ref "./running" >}})

Then we need to set up the test infrastructure following the instructions form [here](https://owncloud.dev/clients/web/testing/)

Now we can run the tests. The tests will take several configuration variables which can be found [here](https://owncloud.dev/clients/web/testing/#available-settings-to-be-set-by-environment-variables). Without configuration, most of the defaults will work. We just need make sure to set these values through env variable.

``` bash
export WEB_PATH=<path to ownCloud Web directory>
export OCIS_SKELETON_DIR=<path to the skeleton directory>
export WEB_UI_CONFIG=<path to the config.json file used by web>
```

While running oCIS we should always use a configuration file for ownCloud Web because our tests will read this file and sometimes even change it which cannot be done if you use environment variables or the default values.

With all this in place we can just run the tests with a simple make command.
First go to the Hello repository
```bash
cd <path to hello>
```
Then run

```bash
make test-acceptance-webui
```

To run just one feature you can run
```bash
make test-acceptance-webui <path-to-feature file>:<line-number>
```
