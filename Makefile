ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin

help: _help_

_help_:
	@echo make build - build go programs in to the bin folder


build:
	cd $(BIN_DIR) && go build $(ROOT_DIR)/main.go



