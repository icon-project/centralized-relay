VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
DIRTY := $(shell git status --porcelain | wc -l | xargs)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

all: lint install

###############################################################################
# Build / Install
###############################################################################

ldflags = -X github.com/icon-project/centralized-relay/cmd.Version=$(VERSION) \
					-X github.com/icon-project/centralized-relay.Commit=$(COMMIT) \
					-X github.com/icon-project/centralized-relay.Dirty=$(DIRTY) \
					-s -w

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -mod=readonly -trimpath -ldflags '$(ldflags)'

build: go.sum
	@echo "building centralized-relay binary..."
	ifeq ($(OS),Windows_NT)
		@go build $(BUILD_FLAGS) -o build/centralized-relay main.go
	else
		@go build -tags=muslc $(BUILD_FLAGS) -o build/centralized-relay main.go
	endif

build-docker:
	@echo "building centralized docker image..."
	docker build . -t centralized-relay

install: go.sum
	@echo "--> Installing centralized-relay binary to $(GOBIN)"
	@go build $(BUILD_FLAGS) -o $(GOBIN)/centralized-relay main.go

install-dev: go.sum
	@echo "installing centralized-relay binary..."
	@go build -mod=readonly -ldflags '$(ldflags)' -o $(GOBIN)/centralized-relay main.go

e2e-test:
	@go test -v ./test/e2e -testify.m TestE2E_all -timeout 30m

test-all:
	@go test -v ./...

PACKAGE_NAME          := github.com/icon-project/centralized-relay
GOLANG_CROSS_VERSION  ?= v1.23.3
LIBWASM_VERSION 		 ?= v2.1.0

SYSROOT_DIR     ?= sysroots
SYSROOT_ARCHIVE ?= sysroots.tar.bz2


.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--env LIBWASM_VERSION=$(LIBWASM_VERSION) \
		--env COSIGN_PASSWORD \
		--env COSIGN_PRIVATE_KEY \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --auto-snapshot

.PHONY: release
release:
	docker run \
		--rm \
		--env GITHUB_TOKEN \
		--env LIBWASM_VERSION=$(LIBWASM_VERSION) \
		--env COSIGN_PWD \
		--env COSIGN_KEY \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean
