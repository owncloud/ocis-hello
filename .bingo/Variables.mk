# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.8. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.8.0
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.8.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.8.0 "github.com/bwplotka/bingo"

CALENS := $(GOBIN)/calens-v0.3.0
$(CALENS): $(BINGO_DIR)/calens.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/calens-v0.3.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=calens.mod -o=$(GOBIN)/calens-v0.3.0 "github.com/restic/calens"

FILEB0X := $(GOBIN)/fileb0x-v1.1.4
$(FILEB0X): $(BINGO_DIR)/fileb0x.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/fileb0x-v1.1.4"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=fileb0x.mod -o=$(GOBIN)/fileb0x-v1.1.4 "github.com/UnnoTed/fileb0x"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v1.52.2
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v1.52.2"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v1.52.2 "github.com/golangci/golangci-lint/cmd/golangci-lint"

GOVERAGE := $(GOBIN)/goverage-v0.0.0-20180129164344-eec3514a20b5
$(GOVERAGE): $(BINGO_DIR)/goverage.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/goverage-v0.0.0-20180129164344-eec3514a20b5"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=goverage.mod -o=$(GOBIN)/goverage-v0.0.0-20180129164344-eec3514a20b5 "github.com/haya14busa/goverage"

GOX := $(GOBIN)/gox-v1.0.1
$(GOX): $(BINGO_DIR)/gox.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/gox-v1.0.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=gox.mod -o=$(GOBIN)/gox-v1.0.1 "github.com/mitchellh/gox"

HUGO := $(GOBIN)/hugo-v0.112.5
$(HUGO): $(BINGO_DIR)/hugo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/hugo-v0.112.5"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=hugo.mod -o=$(GOBIN)/hugo-v0.112.5 "github.com/gohugoio/hugo"

MUTAGEN := $(GOBIN)/mutagen-v0.11.8
$(MUTAGEN): $(BINGO_DIR)/mutagen.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mutagen-v0.11.8"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=mutagen.mod -o=$(GOBIN)/mutagen-v0.11.8 "github.com/mutagen-io/mutagen/cmd/mutagen"

PROTOC_GEN_DOC := $(GOBIN)/protoc-gen-doc-v1.5.1
$(PROTOC_GEN_DOC): $(BINGO_DIR)/protoc-gen-doc.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-doc-v1.5.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-doc.mod -o=$(GOBIN)/protoc-gen-doc-v1.5.1 "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"

PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc-v1.3.0
$(PROTOC_GEN_GO_GRPC): $(BINGO_DIR)/protoc-gen-go-grpc.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-grpc-v1.3.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-go-grpc.mod -o=$(GOBIN)/protoc-gen-go-grpc-v1.3.0 "google.golang.org/grpc/cmd/protoc-gen-go-grpc"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.30.0
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.30.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.30.0 "google.golang.org/protobuf/cmd/protoc-gen-go"

PROTOC_GEN_GRPC_GATEWAY := $(GOBIN)/protoc-gen-grpc-gateway-v2.15.2
$(PROTOC_GEN_GRPC_GATEWAY): $(BINGO_DIR)/protoc-gen-grpc-gateway.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-grpc-gateway-v2.15.2"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-grpc-gateway.mod -o=$(GOBIN)/protoc-gen-grpc-gateway-v2.15.2 "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"

PROTOC_GEN_MICRO := $(GOBIN)/protoc-gen-micro-v1.0.0
$(PROTOC_GEN_MICRO): $(BINGO_DIR)/protoc-gen-micro.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-micro-v1.0.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-micro.mod -o=$(GOBIN)/protoc-gen-micro-v1.0.0 "github.com/go-micro/generator/cmd/protoc-gen-micro"

PROTOC_GEN_MICROWEB := $(GOBIN)/protoc-gen-microweb-v0.0.0-20220808092353-b5d6c3960e19
$(PROTOC_GEN_MICROWEB): $(BINGO_DIR)/protoc-gen-microweb.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-microweb-v0.0.0-20220808092353-b5d6c3960e19"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-microweb.mod -o=$(GOBIN)/protoc-gen-microweb-v0.0.0-20220808092353-b5d6c3960e19 "github.com/owncloud/protoc-gen-microweb"

PROTOC_GEN_OPENAPIV2 := $(GOBIN)/protoc-gen-openapiv2-v2.15.2
$(PROTOC_GEN_OPENAPIV2): $(BINGO_DIR)/protoc-gen-openapiv2.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-openapiv2-v2.15.2"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-openapiv2.mod -o=$(GOBIN)/protoc-gen-openapiv2-v2.15.2 "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"

PROTOC_GO_INJECT_TAG := $(GOBIN)/protoc-go-inject-tag-v1.4.0
$(PROTOC_GO_INJECT_TAG): $(BINGO_DIR)/protoc-go-inject-tag.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-go-inject-tag-v1.4.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-go-inject-tag.mod -o=$(GOBIN)/protoc-go-inject-tag-v1.4.0 "github.com/favadi/protoc-go-inject-tag"

