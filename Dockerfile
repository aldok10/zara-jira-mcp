FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/bin/zara-jira-mcp ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/zara-jira-mcp /usr/local/bin/zara-jira-mcp

VOLUME /data
ENV PM_MEMORY_DB_PATH=/data/pm_memory.db

ENTRYPOINT ["zara-jira-mcp"]
