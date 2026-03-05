BINARY_NAME=bin/kryvora-node

.PHONY: all build clean test run

all: build

build:
	mkdir -p bin
	go build -o $(BINARY_NAME) ./cmd/kryvora-node

clean:
	rm -rf bin/

test:
	go test -v ./...

run: build
	./$(BINARY_NAME) --config config.example.yaml
