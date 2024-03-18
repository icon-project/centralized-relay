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
					-X github.com/icon-project/centralized-relay.Dirty=$(DIRTY)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building centralized-relay binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/centralized-relayer.exe main.go
else
	@echo "building centralized-relayer binary..."
	@go build  $(BUILD_FLAGS) -o build/centralized-relay main.go
endif

build-docker:
	@echo "building centralized docker image..."
	docker build . -t centralized-relay

install: go.sum
	@echo "installing centralized-relay binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/centralized-relay main.go

PACKAGE_NAME          := github.com/icon-project/centralized-relay
GOLANG_CROSS_VERSION  ?= v1.22.0

SYSROOT_DIR     ?= sysroots
SYSROOT_ARCHIVE ?= sysroots.tar.bz2


.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		-e CGO_ENABLED=0 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist