SHELL := /bin/bash
.SHELLFLAGS := -c

ORG_NAME ?= kegr.io
PROJECT_NAME ?= storage_controller

GOPATH := $(shell pwd)
BIN_DIR ?= $(shell pwd)/target
ENV_VARS ?= GOPATH=$(GOPATH) PATH=$(shell echo $$PATH):$(GOPATH)/bin
SRC_DIR ?= ./src
GO_DEPS ?= $(shell cat .godep)

PROTO_PACKAGE ?= protobuf
PROTO_PATH ?= ./$(ORG_NAME)/$(PROTO_PACKAGE)
PROTO_DIRS := $(shell cd $(SRC_DIR)/$(PROTO_PATH); find . -name "*.proto" -type f -print0 | xargs -0 -n1 dirname | sort | uniq)

.PHONY: clean build run dep proto

build: proto
	$(INFO) "building"
	@$(ENV_VARS) go build -o "$(BIN_DIR)/$(PROJECT_NAME)" $(SRC_DIR)/$(ORG_NAME)/$(PROJECT_NAME)

run: build
	$(INFO) "running"	
	@cd $(BIN_DIR); ./$(PROJECT_NAME)

clean:
	$(INFO) "deleting build directory $(BIN_DIR)"
	@rm -rf $(BIN_DIR)
	@mkdir $(BIN_DIR)
	@mkdir $(BIN_DIR)/www

dep:
	$(INFO) "installing dependencies"
	$(foreach \
		dep, \
		$(GO_DEPS), \
		$(shell $(ENV_VARS) go get -u $(dep)) \
	)

proto:
	$(INFO) "deleting old generated proto files"
	@$(shell cd $(SRC_DIR)/$(PROTO_PATH); find . -name "*.pb.go" -type f -exec rm {} \;)
	$(INFO) "generating proto files"
	$(foreach \
		protodir, \
		$(PROTO_DIRS), \
		$(shell cd $(SRC_DIR); protoc -I=$(PROTO_PATH) --go_out=plugins=grpc:. $(PROTO_PATH)/$(protodir)/*.proto) \
	)

YELLOW := "\e[1;33m"
NC := "\e[0m"
INFO := @bash -c '\
  printf $(YELLOW); \
  echo "=> $$1"; \
  printf $(NC)' SOME_VALUE
