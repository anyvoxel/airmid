# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build' target
#   make lint         - code analysis
#   make test         - run unit test (or plus integration test)
#   make clean        - clean up targets
#
# Not included but recommended targets:
#   make e2e-test
#
# The makefile is also responsible to populate project version information.
#

#
# Tweak the variables based on your project.
#

# This repo's root import path (under GOPATH).
ROOT := github.com/anyvoxel/airmid/ioc

# Module name.
NAME := ioc

#
# These variables should not need tweaking.
#


#
# Define all targets. At least the following commands are required:
#

# All targets.
.PHONY: lint test clean mock

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) --version
	@$(GOLANGCI_LINT) run -v --config ../.golangci.yaml


test:
	# 1. Now we can use -race flag, the 'nosplit stack overflow' error will not happened at >= go1.20
	#   NOTE: see https://github.com/golang/go/issues/54291 for more detail
	# 2. add -ldflags to prevent 'permission denied' in macos, see https://github.com/agiledragon/gomonkey/issues/70 for more detail.
	@go test -v -race -ldflags="-extldflags="-Wl,-segprot,__TEXT,rwx,rx"" -coverpkg=./... -coverprofile=coverage.out -gcflags="all=-N -l" ./...
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'


.PHONY: clean
clean:
	@-rm -vrf ${OUTPUT_DIR} output coverage.out


.PHONY: mock
mock:
	mockgen -source=props/types.go -destination=props/mocks/types.go -package=mocks
