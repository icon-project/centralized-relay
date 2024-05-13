VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
DIRTY := $(shell git status --porcelain | wc -l | xargs)

GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Define the URLs for the library files
LIBWASMV_VERS := "v2.0.0"
LIBWASMMV_URL := https://github.com/CosmWasm/wasmvm/releases/download/$(LIBWASMV_VERS)



all: lint install

###############################################################################
# Build / Install
###############################################################################

ldflags := -s -w -X github.com/icon-project/centralized-relay/cmd.Version=$(VERSION) \
					-X github.com/icon-project/centralized-relay.Commit=$(COMMIT) \
					-X github.com/icon-project/centralized-relay.Dirty=$(DIRTY) \
					-linkmode=external

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

# Default build flags for each OS
ifeq ($(GOOS),linux)
	BUILD_TAGS := -tags muslc,netgo
	ifeq ($(GOARCH),amd64)
		BUILD_FLAGS := -ldflags '$(ldflags) -extldflags "-Wl,-z,muldefs -lm"'
		BUILD_ENV := CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc CXX=x86_64-linux-gnu-g++
	else ifeq ($(GOARCH),arm64)
		BUILD_FLAGS := -ldflags '$(ldflags) -extldflags "-static"'
		BUILD_ENV := CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc
	endif
else ifeq ($(GOOS),darwin)
	BUILD_TAGS := -tags static_wasm
	ifeq ($(GOARCH),amd64)
		BUILD_FLAGS := -ldflags '$(ldflags) -extldflags "-Wl,-z,muldefs -lm"'
		BUILD_ENV := CGO_ENABLED=1 CC=o64-clang CGO_LDFLAGS=-L/lib

	else ifeq ($(GOARCH),arm64)
		BUILD_FLAGS := -ldflags '$(ldflags) -extldflags "-static"'
		BUILD_ENV := CGO_ENABLED=1 CC=o64-clang CGO_LDFLAGS=-L/lib
	endif
else
	BUILD_FLAGS := -ldflags '$(ldflags)'
	BUILD_ENV := CGO_ENABLED=1
endif

install_libwasmvm:
	@echo "Install libwasmvm static library..."
	@if [ "$(GOOS)" = "linux" ] && [ "$(GOARCH)" = "amd64" ]; then \
		sudo wget -q $(LIBWASMMV_URL)/libwasmvm_muslc.x86_64.a -O /usr/lib/x86_64-linux-gnu/libwasmvm_muslc.a; \
	elif [ "$(GOOS)" = "darwin" ] && [ "$(GOARCH)" = "amd64" ]; then \
		sudo wget -q $(LIBWASMMV_URL)/libwasmvm_muslc.darwin.a -O /usr/local/lib/libwasmvm_muslc.a; \
	elif [ "$(GOOS)" = "linux" ] && [ "$(GOARCH)" = "arm64" ]; then \
		sudo wget -q $(LIBWASMMV_URL)/libwasmvm_muslc.aarch64.a -O /usr/lib/aarch64-linux-gnu/libwasmvm_muslc.a; \
	elif [ "$(GOOS)" = "darwin" ] && [ "$(GOARCH)" = "arm64" ]; then \
		sudo wget -q $(LIBWASMMV_URL)/libwasmvmstatic_darwin.a -O /lib/libwasmvmstatic_darwin.a; \
	else \
		echo "Unsupported GOOS=$(GOOS) GOARCH=$(GOARCH) combination for libwasmvm installation."; \
		exit 1; \
	fi

build: go.sum install_libwasmvm
ifeq ($(GOOS),windows)
	@echo "building centralized-relay binary for Windows..."
	@$(BUILD_ENV) go build -mod=readonly -trimpath $(BUILD_FLAGS) $(BUILD_TAGS) -o build/centralized-relay.exe main.go
else ifeq ($(GOOS),linux)
	@echo "building centralized-relay binary for Linux..."
	@if [ "$(GOARCH)" = "amd64" ]; then \
		$(BUILD_ENV) go build -mod=readonly -trimpath $(BUILD_FLAGS) $(BUILD_TAGS) -o build/centralized-relay main.go; \
	elif [ "$(GOARCH)" = "arm64" ]; then \
		$(BUILD_ENV) go build -mod=readonly -trimpath $(BUILD_FLAGS) $(BUILD_TAGS) -o build/centralized-relay main.go; \
	else \
		echo "Unsupported GOARCH=$(GOARCH) for $(GOOS)"; \
		exit 1; \
	fi
else ifeq ($(GOOS),darwin)
	@echo "building centralized-relay binary for macOS..."
	@if [ "$(GOARCH)" = "amd64" ]; then \
		$(BUILD_ENV) go build -mod=readonly -trimpath $(BUILD_FLAGS) $(BUILD_TAGS) -o build/centralized-relay main.go; \
	elif [ "$(GOARCH)" = "arm64" ]; then \
		$(BUILD_ENV) go build -mod=readonly -trimpath $(BUILD_FLAGS) $(BUILD_TAGS) -o build/centralized-relay main.go; \
	else \
		echo "Unsupported GOARCH=$(GOARCH) for $(GOOS)"; \
		exit 1; \
	fi
else
	@echo "Unsupported GOOS=$(GOOS)"
	exit 1
endif

build-docker:
	@echo "building centralized docker image..."
	docker build . -t centralized-relay

install: go.sum install_libwasmvm build
	@echo "installing centralized-relay binary..."
	@mv ./build/centralized-relay $(GOBIN)/centralized-relay
	@rm -rf ./build

install-dev: go.sum
	@echo "installing centralized-relay binary..."
	@go build -mod=readonly -ldflags '$(ldflags)' -o $(GOBIN)/centralized-relay main.go

e2e-test:
	@go test -v ./test/e2e -testify.m TestE2E_all

PACKAGE_NAME          := github.com/icon-project/centralized-relay
GOLANG_CROSS_VERSION  ?= v1.22.1
COSMWASM_VERSION      ?= v2.0.0

SYSROOT_DIR     ?= sysroots
SYSROOT_ARCHIVE ?= sysroots.tar.bz2


.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish

.PHONY: release
release:
	docker run \
		--rm \
		--env GITHUB_TOKEN \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean