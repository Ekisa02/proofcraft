BINARY=proofcraft

.PHONY: build run test

build:
	go build -o $(BINARY) main.go

run: build
	./$(BINARY)

test:
	go test ./...
