SHELL := bash
NAME := hello

.PHONY: test-acceptance-webui
test-acceptance-webui:
	yarn run test:e2e $(FEATURE_PATH)


############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # supress `command not found warnings` for non go targets in CI
include .bingo/Variables.mk
endif

############ go tooling ############
include .make/go.mk

############ release ############
include .make/release.mk

############ docs generate ############
include .make/docs.mk

############ l10n ############
include .make/l10n.mk

.PHONY: docs-generate
docs-generate: config-docs-generate \
			   grpc-docs-generate

############ generate ############
include .make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: protobuf # CI runs ci-node-generate automatically before this target
	@go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate: yarn-build

.PHONY: yarn-build
yarn-build: node_modules
	yarn lint
	yarn test
	yarn build

.PHONY: node_modules
node_modules:
	yarn install --frozen-lockfile

############ protobuf ############
PROTO_VERSION := v0
PROTO_SRC := pkg/proto/$(PROTO_VERSION)

include .make/protobuf.mk

.PHONY: protobuf
protobuf: $(PROTO_SRC)/${NAME}.pb.go \
		  $(PROTO_SRC)/${NAME}.pb.micro.go \
		  docs/extensions/${NAME}/grpc.md


.PHONY: bingo-update
bingo-update: $(BINGO)
	$(BINGO) get -u
