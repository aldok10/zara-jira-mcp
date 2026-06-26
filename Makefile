.PHONY: build install run test clean lint docker-build docker-run

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

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

docker-build:
	docker build -t $(BINARY):$(VERSION) .

docker-run:
	docker compose up --build
