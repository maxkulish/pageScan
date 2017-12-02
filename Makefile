ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin

help: _help_

_help_:
	@echo make build - build go programs in to the ./bin folder


build_mac:
	cd $(ROOT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o ./bin/pageScanMac

build_linux:
	cd $(ROOT_DIR)
	GOOS=linux GOARCH=amd64 go build -o ./bin/pageScanLinux

build:
	make build_mac
	make build_linux




