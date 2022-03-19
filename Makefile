SHELL:=/bin/sh
.PHONY: build build_server \
		test run fmt vet clean \
		mod_update vendor_from_mod vendor_clean

export GO111MODULE=on

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}bin
GO_PATH := $(shell go env | grep GOPATH | awk -F '"' '{print $$2}')
HTTPSERVER_TEST_PATH := build/test

# Image Name
IMAGE_NAME?=mawenzhong/gwork

# Version
RELEASE?=v1.0.0

# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

# Build Flags
GO_LD_FLAGS= "-s -w"

# Cgo is disabled by default
ENABLE_CGO= CGO_ENABLED=0

# Check Go build tags, the tags are from command line of make
ifdef GOTAGS
  GO_BUILD_TAGS= -tags ${GOTAGS}
  # Must enable Cgo when wasmhost is included
  ifeq ($(findstring wasmhost,${GOTAGS}), wasmhost)
	ENABLE_CGO= CGO_ENABLED=1
  endif
endif

# When build binaries for docker, we put the binaries to another folder to avoid
# overwriting existing build result, or Mac/Windows user will have to do a rebuild
# after build the docker image, which is Linux only currently.
ifdef DOCKER
  RELEASE_DIR= ${MKFILE_DIR}build/bin
endif

# Targets
TARGET_SERVER=${RELEASE_DIR}/server

# Rules

build_server:
	@echo "build server"
	cd ${MKFILE_DIR} && \
	${ENABLE_CGO} go build ${GO_BUILD_TAGS} -v -trimpath -ldflags ${GO_LD_FLAGS} \
	-o ${TARGET_SERVER} ${MKFILE_DIR}cmd/server.go

build_linux_server:
	@echo "build server"
	cd ${MKFILE_DIR} && \
	${ENABLE_CGO} GOOS=linux go build ${GO_BUILD_TAGS} -v -trimpath -ldflags ${GO_LD_FLAGS} \
	-o ${TARGET_SERVER} ${MKFILE_DIR}cmd/server.go

dev_build_server:
	@echo "build dev server"
	cd ${MKFILE_DIR} && \
	go build -v -race -ldflags ${GO_LD_FLAGS} \
	-o ${TARGET_SERVER} ${MKFILE_DIR}cmd/server.go

test:
	cd ${MKFILE_DIR}
	go mod tidy
	git diff --exit-code go.mod go.sum
	go mod verify
	go test -v ./... ${TEST_FLAGS}

httpserver_test: build
	{ \
	set -e ;\
	cd ${HTTPSERVER_TEST_PATH} ;\
	./httpserver_test.sh ;\
    }

clean:
	rm -rf ${RELEASE_DIR}
	rm -rf ${MKFILE_DIR}build/cache
	rm -rf ${MKFILE_DIR}build/bin

start_server:
	${TARGET_SERVER}

run: dev_build_server start_server

build_docker: build_linux_server
	docker build -t ${IMAGE_NAME}:${RELEASE} -f ./Dockerfile .
	docker tag ${IMAGE_NAME}:${RELEASE} ${IMAGE_NAME}:latest

	docker login --username=mawenzhong
	docker push ${IMAGE_NAME}:latest
	docker image prune -f

fmt:
	cd ${MKFILE_DIR} && go fmt ./...

vet:
	cd ${MKFILE_DIR} && go vet ./...

vendor_from_mod:
	cd ${MKFILE_DIR} && go mod vendor

vendor_clean:
	rm -rf ${MKFILE_DIR}vendor

mod_update:
	cd ${MKFILE_DIR} && go get -u
