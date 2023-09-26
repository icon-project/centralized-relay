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
	@echo "building centralized rly binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/centralized-rly.exe main.go
else
	@echo "building centralized rly binary..."
	@go build  $(BUILD_FLAGS) -o build/centralized-rly main.go
endif


install: go.sum
	@echo "installing centralized rly binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/centralized-rly main.go