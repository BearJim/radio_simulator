BIN_PATH = bin
SIMULATOR = simulator
SIMCTL = simctl

.PHONY: $(SIMULATOR) clean

.DEFAULT_GOAL: SIMULATOR

all: $(SIMULATOR) $(SIMCTL)

$(SIMULATOR): cmd/$(SIMULATOR)/main.go
	go build -o $(BIN_PATH)/$@ $^

$(SIMCTL): cmd/$(SIMCTL)/main.go
	go build -o $(BIN_PATH)/$@ $^

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/api/api.proto

clean:
	rm -rf $(BIN_PATH)/
