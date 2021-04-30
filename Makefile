BIN_PATH = bin
SIMULATOR = simulator
SIMCTL = simctl

VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH = $(shell git log --pretty='%H' -1 | cut -c1-8)

LDFLAGS = -X github.com/free5gc/version.VERSION=$(VERSION) \
          -X github.com/free5gc/version.BUILD_TIME=$(BUILD_TIME) \
          -X github.com/free5gc/version.COMMIT_HASH=$(COMMIT_HASH)

.PHONY: $(SIMULATOR) clean

.DEFAULT_GOAL: SIMULATOR

all: $(SIMULATOR) $(SIMCTL)

$(SIMULATOR): cmd/$(SIMULATOR)/main.go
	go build -ldflags "$(LDFLAGS)" -o $(BIN_PATH)/$@ $^

$(SIMCTL): cmd/$(SIMCTL)/main.go
	go build -o $(BIN_PATH)/$@ $^

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/api/api.proto

clean:
	rm -rf $(BIN_PATH)/
