.PHONY: all build export run
GOOS:=$(shell go env GOOS)
GOARCH:=amd64
EXECUTABLE:=funservice
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export PROJECTNAME:=$(shell basename "$(ROOT_DIR)")

# Git parameters
VERSION := $(shell git describe --tags --always --dirty="-dev")
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

ifdef FUNSERVICE_DIR
else
	FUNSERVICE_DIR:=../../funservice/build/$(GOOS)
endif

PATH := $(PATH):$(FUNSERVICE_DIR)

GOCMD:=go
GOTEST:=$(GOCMD) test
GOTOOL:=$(GOCMD) tool
GOFORMAT:=$(GOCMD) fmt
GOLIST:=$(GOCMD) list

export GO111MODULE=on
export GOPROXY=https://goproxy.io

all: build run
build:
	$(EXECUTABLE) build $(PROJECTNAME) $(VERSION)_$(BUILDTIME) $(ROOT_DIR) $(FUNSERVICE_DIR)/plugins/$(subst Plugin,.plugin,$(PROJECTNAME))

export:
	$(EXECUTABLE) build $(PROJECTNAME) $(VERSION)_$(BUILDTIME) $(ROOT_DIR) $(ROOT_DIR)/build/$(subst Plugin,_$(VERSION)_$(BUILDTIME).plugin,$(PROJECTNAME))
	
run:
	$(EXECUTABLE) serve -v --plugin $(subst Plugin,,$(PROJECTNAME))

setup:
	$(EXECUTABLE) init

test:
	$(GOTEST) ./...

fmt:
	$(GOFORMAT) $$($(GOLIST) ./...)
	git diff --exit-code

stop-docker:
	docker-compose -f ./docker-compose.yaml down
	docker volume rm $(subst Plugin,plugin,$(PROJECTNAME))_fundata

build-docker:
	@docker-compose up -d
	@docker-compose exec service /bin/bash -c "make build;/bin/bash"
	