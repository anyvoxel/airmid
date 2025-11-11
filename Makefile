# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'all' target
#   make lint         - code analysis
#   make test         - run unit test (or plus integration test)
#   make clean        - clean up targets
#   make mock         - generate mock files
#
#
# The makefile is also responsible to populate project version information.
#


#
# These variables should not need tweaking.
#

# It's necessary to set this because some environments don't link sh -> bash.
export SHELL := /bin/bash

# It's necessary to set the errexit flags for the bash shell.
export SHELLOPTS := errexit

# Project output directory.
export OUTPUT_DIR := ./bin

# Build directory.
export BUILD_DIR := ./build

# Current version of the project.
export VERSION      ?= $(shell git describe --tags --always --dirty)
export BRANCH       ?= $(shell git branch | grep \* | cut -d ' ' -f2)
export GITCOMMIT    ?= $(shell git rev-parse HEAD)
export GITTREESTATE ?= $(if $(shell git status --porcelain),dirty,clean)
export BUILDDATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export appVersion   ?= $(VERSION)

# Available cpus for compiling, please refer to https://major.io/2019/04/05/inspecting-openshift-cgroups-from-inside-the-pod/ for more information.
export CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)

# Track code version with Docker Label.
export DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
export BIN_DIR := $(GOPATH)/bin
export GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

# Default golang flags used in build and test
# -count: run each test and benchmark 1 times. Set this flag to disable test cache
export GOFLAGS ?= -count=1


# NOTE：we cannot use `go list` to retrieve all modules, because go.work won't submit to 
# version control system.
# GOMODS ?= $(shell go list -f '{{.Dir}}' -m | sed "s|^$$(pwd)/||")
# NOTE: we can use `find . -name "go.mod" -exec dirname {} \; | xargs realpath --relative-to=.` to 
# automatically find all modules.
GOMODS ?= anvil app ioc

# #
# # Define all targets. At least the following commands are required:
# #

.DEFAULT_GOAL := all


.PHONY: help
help:  ## list all available targets
	@echo "available targets:"
	@echo "=========="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' | \
	sort


.PHONY: all
all: test lint clean  ## run test、lint、clean target


.PHONY: test
test:  ## 
	# 1. remove -race flag to prevent 'nosplit stack overflow' error, see https://github.com/golang/go/issues/54291 for more detail
	# 	NOTE: this was fixed by 1.20 release
	# 2. add -ldflags to prevent 'permission denied' in macos, see https://github.com/agiledragon/gomonkey/issues/70 for more detail.
	@for gomod in $(GOMODS); do \
	    (cd $$gomod && \
			go test -v -race -ldflags="-extldflags="-Wl,-segprot,__TEXT,rwx,rx"" -coverpkg=./... -coverprofile=coverage.out -gcflags="all=-N -l" ./... && \
			go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'); \
	done


.PHONY: lint
lint: $(GOLANGCI_LINT) ## 
	$(GOLANGCI_LINT) --version
	@for gomod in $(GOMODS); do \
		$(GOLANGCI_LINT) run -v --config ./.golangci.yaml $$gomod/... ; \
	done

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(BIN_DIR) v2.4.0


.PHONY: clean
clean:  ##
	@find . -name "coverage.out*" -exec rm -vrf {} \;
	@for gomod in $(GOMODS); do \
	    (cd $$gomod && rm -vrf ${OUTPUT_DIR} output coverage.out); \
	done


MOCKGEN := $(BIN_DIR)/mockgen
.PHONY: mock
mock: $(MOCKGEN)  ## 
	mockgen -source=ioc/props/types.go -destination=ioc/props/mocks/types.go -package=mocks
	mockgen -source=app/xapp.go -destination=app/xapp_mocks_test.go -package=app
	mockgen -source=app/locator.go -destination=app/locator_mocks_test.go -package=app

$(MOCKGEN):
	go install go.uber.org/mock/mockgen@latest


addheaders: ## add license header to files
	@command -v addlicense > /dev/null || go install -v github.com/google/addlicense@v0.0.0-20210428195630-6d92264d7170
	@addlicense -c "The anyvoxel Authors" -l mit .


PROTOCGO := $(BIN_DIR)/protoc-gen-go
PROTOCGRPC := $(BIN_DIR)/protoc-gen-go-grpc
PROTOCGATEWAY := $(BIN_DIR)/protoc-gen-grpc-gateway
.PHONY: proto
proto: $(PROTOCGO) $(PROTOCGRPC) $(PROTOCGATEWAY)  ## 
	@rm -rf ./pb
	@./proto/generate.sh

$(PROTOCGO):
	go install -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0

$(PROTOCGRPC):
	go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

$(PROTOCGATEWAY):
	go install -v github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	go install -v github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
