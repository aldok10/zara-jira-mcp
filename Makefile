.PHONY: build install run test test-cover lint clean tidy docker-build docker-run

BINARY := zara-jira-mcp
VERSION := 0.4.0

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/server

install: build
	cp bin/$(BINARY) $(HOME)/.local/bin/$(BINARY)

run:
	go run ./cmd/server

test:
	go test ./... -count=1

test-cover:
	go test ./application/tools/ -cover -count=1

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
