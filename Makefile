SHELL := bash
NAME := ocis-hello
IMPORT := github.com/owncloud/$(NAME)
BIN := bin
DIST := dist

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(NAME).exe
else
	EXECUTABLE := $(NAME)
endif

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f -not -path "./node_modules/*")
GENERATE ?= $(IMPORT)/pkg/assets

TAGS ?=

ifndef OUTPUT
	ifneq ($(DRONE_TAG),)
		OUTPUT ?= $(subst v,,$(DRONE_TAG))
	else
		OUTPUT ?= testing
	endif
endif

ifndef VERSION
	ifneq ($(DRONE_TAG),)
		VERSION ?= $(subst v,,$(DRONE_TAG))
	else
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

LDFLAGS += -s -w -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	go clean -i ./...
	rm -rf $(BIN) $(DIST) pkg/assets/embed.go

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: staticcheck
staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -tags '$(TAGS)' $(PACKAGES)

.PHONY: lint
lint:
	for PKG in $(PACKAGES); do go run golang.org/x/lint/golint -set_exit_status $$PKG || exit 1; done;

.PHONY: generate
generate: protobuf
	go generate $(GENERATE)

.PHONY: changelog
changelog:
	go run github.com/restic/calens >| CHANGELOG.md

.PHONY: test
test:
	go run github.com/haya14busa/goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: install
install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$(NAME)
	go install -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$(NAME)-cli

.PHONY: build
build: $(BIN)/$(EXECUTABLE) $(BIN)/$(EXECUTABLE)-cli

$(BIN)/$(EXECUTABLE): $(SOURCES)
	go build -i -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-cli: $(SOURCES)
	go build -i -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)-cli

.PHONY: release
release: release-dirs release-linux release-windows release-darwin release-copy release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-linux
release-linux: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'linux' -arch 'amd64 386 arm64 arm' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-windows
release-windows: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'windows' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-darwin
release-darwin: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -os 'darwin' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

.PHONY: release-finish
release-finish: release-copy release-check

.PHONY: docs
docs:
	cd docs; hugo

.PHONY: watch
watch:
	go run github.com/cespare/reflex -c reflex.conf

pkg/proto/v0/hello.pb.go: pkg/proto/v0/hello.proto
	protoc \
		-I=third_party/ \
		-I=pkg/proto/v0/ \
		--go_out=logtostderr=true:pkg/proto/v0 hello.proto

pkg/proto/v0/hello.pb.micro.go: pkg/proto/v0/hello.proto
	protoc \
		-I=third_party/ \
		-I=pkg/proto/v0/ \
		--micro_out=logtostderr=true:pkg/proto/v0 hello.proto

pkg/proto/v0/hello.pb.web.go: pkg/proto/v0/hello.proto
	protoc \
		-I=third_party/ \
		-I=pkg/proto/v0/ \
		--microweb_out=logtostderr=true:pkg/proto/v0 hello.proto

pkg/proto/v0/hello.swagger.json: pkg/proto/v0/hello.proto
	protoc \
		-I=third_party/ \
		-I=pkg/proto/v0/ \
		--swagger_out=logtostderr=true:pkg/proto/v0 hello.proto

.PHONY: protobuf
protobuf:  $(GOPATH)/bin/protoc-gen-go $(GOPATH)/bin/protoc-gen-micro $(GOPATH)/bin/protoc-gen-microweb $(GOPATH)/bin/protoc-gen-swagger pkg/proto/v0/hello.pb.go pkg/proto/v0/hello.pb.micro.go pkg/proto/v0/hello.pb.web.go pkg/proto/v0/hello.swagger.json

$(GOPATH)/bin/protoc-gen-go:
	GO111MODULE=off go get -v github.com/golang/protobuf/protoc-gen-go

$(GOPATH)/bin/protoc-gen-micro:
	GO111MODULE=off go get -v github.com/micro/protoc-gen-micro

$(GOPATH)/bin/protoc-gen-microweb:
	GO111MODULE=off go get -v github.com/webhippie/protoc-gen-microweb

$(GOPATH)/bin/protoc-gen-swagger:
	GO111MODULE=off go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
