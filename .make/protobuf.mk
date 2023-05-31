.PHONY: $(PROTO_SRC)/${NAME}.pb.go
$(PROTO_SRC)/${NAME}.pb.go: $(BUF) $(PROTOC_GEN_GO)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.pb.go"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=./third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-go=$(PROTOC_GEN_GO) \
		--go_out=$(PROTO_SRC) --go_opt=paths=source_relative \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.pb.micro.go
$(PROTO_SRC)/${NAME}.pb.micro.go: $(BUF) $(PROTOC_GEN_MICRO)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.pb.micro.go"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=./third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-micro=$(PROTOC_GEN_MICRO) \
		--micro_out=$(PROTO_SRC) --micro_opt=paths=source_relative \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: ./docs/extensions/$(NAME)/grpc.md
./docs/extensions/$(NAME)/grpc.md: $(BUF) $(PROTOC_GEN_DOC)
	@echo "$(NAME): ./docs/extensions/$(NAME)/grpc.md"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=./third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-doc=$(PROTOC_GEN_DOC) \
		--doc_opt=./templates/GRPC.tmpl,grpc.md \
		--doc_out=./docs/extensions/$(NAME)/ \
		$(PROTO_SRC)/${NAME}.proto
