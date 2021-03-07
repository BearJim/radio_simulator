SIMULATOR = simulator

.PHONY: $(SIMULATOR) clean

.DEFAULT_GOAL: SIMULATOR

all: $(SIMULATOR)

$(SIMULATOR): cmd/$(SIMULATOR)/main.go
	go build -o bin/$@ $^

clean:
	rm $(SIMULATOR)
