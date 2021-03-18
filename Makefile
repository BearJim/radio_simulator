SIMULATOR = simulator

.PHONY: $(SIMULATOR) clean

.DEFAULT_GOAL: SIMULATOR

all: $(SIMULATOR)

$(SIMULATOR): cmd/$(SIMULATOR)/main.go
	go build -o bin/$@ $^

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/api/api.proto

clean:
	rm $(SIMULATOR)
