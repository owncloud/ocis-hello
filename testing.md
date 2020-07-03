---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: testing.md
---
This repository provides a general guideline for creating tests for the ocis extensions. The tests can be written in various levels from unit, integration, and end-to-end. It is not essential to write tests on all these levels as it can be redundant in some cases. This repository provides a reference for all levels of tests.

### Unit tests
Unit tests generally live inside *_test.go files in the `/pkg` directory. One such example in this extension is in `/pkg/service/v0/service_test.go`. Similarly the unit test for the protobuf generated code can also be written just like in `/pkg/proto/hello.pb_test.go`. 

### Integration tests
There are mainly 2 types of integration tests, namely HTTP tests, and GRPC tests. These tests mostly live in `/pkg/proto` directory where all the protobuf definitions are specified. The examples for the HTTP integration tests are in `/pkg/proto/hello.pb.web_test.go` whereas the GRPC tests are in `/pkg/proto/hello.pb.micro_test.go`.

### End-to-End tests
For extensions with an UI, we can also write end-to-end tests using the Nightwatch test framework. These tests live in `/ui/tests` directory. We can reuse already existing Gherkin steps from the [phoenix](https://github.com/owncloud/phoenix) tests here.

### Running the tests
#### Unit and integration tests
The unit and integration tests are run using the simple `go test` command. If you wish to run all the tests with the coverage you can just use make command.
```bash
make test
```
You can also run a specific file with the go test command
```bash
go test <path to package or file>
```
#### End-to-End tests
Running end-to-end tests is a bit more complicated than unit and integration tests. First of all we will need a complete ocis setup with ocis-hello running. For that refer to [this guide](https://owncloud.github.io/extensions/ocis_hello/configuration/).

After that, We need to set the proper test environment.
To run the end-to-end tests, first-of-all we need the phoenix repository where all the test infrastructure exists. So, clone the phoenix repository in your system in any location.
```bash
git clone https://github.com/owncloud/phoenix $HOME/phoenix
```

Next we will need to start the selenium server which will control the browser. There is a script in the phoenix repo that starts the selenium server, just run that to start selenium.
```bash
cd $HOME/phoenix
yarn run selenium
```

Now we can run the tests. The tests will take several configuration variables which can be found [here](https://owncloud.github.io/clients/web/testing/#available-settings-to-be-set-by-environment-variables). Without configuration, most of the defaults will work. We just need make sure to set these values through env variable.

``` bash
export PHOENIX_PATH=$HOME/phoenix
export OCIS_SKELETON_DIR=<path to the skeleton directory>
export PHOENIX_CONFIG=<path to the config.json file used by phoenix>
```
The phoenix path should be set to the directory where the phoenix source files are. Our tests use the existing infrastructure from the phoenix directory to run the tests.

The skeleton directory for the webui tests can be found in [the testing app](https://github.com/owncloud/testing/tree/master/data/webUISkeleton). You can just clone that repository in your local machine and point the env variable to the correct path.

While running ocis we should always use a config file for phoenix because our tests will read this file and sometimes even change it which cannot be done if you use env variables or the default values.

With all this in place we can just run the tests with a simple make command.
First go to the ocis-hello repository
```bash
cd <path to ocis-hello>
```
Then Simply run

```bash
make test-acceptance-webui
```

To run just one feature you can run
```bash
make test-acceptance-webui <path-to-feature file>:<line-number>
```
