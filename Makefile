.PHONY: build install run test test-cover lint clean tidy docker-build docker-run

BINARY := zara-jira-mcp
VERSION := 0.4.0

BINARY := zara-jira-mcp
VERSION := 0.4.0
API_DIR := ./apps/api

build:
	cd $(API_DIR) && go build -ldflags="-s -w -X main.version=$(VERSION)" -o ../../bin/$(BINARY) ./cmd/server

install: build
	cp bin/$(BINARY) $(HOME)/.local/bin/$(BINARY)

run:
	cd $(API_DIR) && go run ./cmd/server

test:
	go test ./bootstrap/...
	cd $(API_DIR) && go test ./...

test-cover:
	cd $(API_DIR) && go test ./... -cover -count=1

lint:
	golangci-lint run ./... ./$(API_DIR)/...

clean:
	rm -rf bin/

tidy:
	go mod tidy
	cd ./shared && go mod tidy
	cd ./modules/jira && go mod tidy
	cd ./modules/sprint && go mod tidy
	cd ./modules/notification && go mod tidy
	cd $(API_DIR) && go mod tidy

docker-build:
	docker build -t $(BINARY):$(VERSION) .

docker-run:
	docker compose up --build
