# Project settings
MODULE=github.com/ahmad-masud/Kivi
APP_NAME=Kivi
PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/kv.proto
PROTO_OUT=.

# Tools
PROTOC=protoc
PROTOC_GEN_GO=$(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)

# Go commands
GO=go

.PHONY: all proto server client clean

all: proto

proto:
	@if [ -z "$(PROTOC_GEN_GO)" ] || [ -z "$(PROTOC_GEN_GO_GRPC)" ]; then \
		echo "Installing protoc plugins..."; \
		$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	fi
	$(PROTOC) --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_FILE)

server:
	$(GO) run ./server

client:
	$(GO) run ./client

clean:
	$(GO) clean
	rm -f $(PROTO_DIR)/*.pb.go
