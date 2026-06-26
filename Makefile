.PHONY: build install run test clean

BINARY := zara-jira-mcp
VERSION := 0.1.0

build:
	go build -ldflags="-s -w" -o bin/$(BINARY) ./cmd/server

install: build
	cp bin/$(BINARY) $(HOME)/.local/bin/$(BINARY)

run:
	go run ./cmd/server

test:
	go test ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy
