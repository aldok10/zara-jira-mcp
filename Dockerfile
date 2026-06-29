FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Workspace configuration
COPY go.work go.work.sum ./
# Module go.mod files
COPY apps/api/go.mod apps/api/go.sum apps/api/
COPY shared/go.mod shared/go.sum shared/
COPY modules/jira/go.mod modules/jira/go.sum modules/jira/
COPY modules/sprint/go.mod modules/sprint/go.sum modules/sprint/
COPY modules/notification/go.mod modules/notification/go.sum modules/notification/
# Root module (required by bootstrap/internal deps)
COPY go.mod go.sum ./

RUN go mod download

# Copy all source
COPY . .

# Build from apps/api
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/bin/zara-jira-mcp ./apps/api/cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/zara-jira-mcp /usr/local/bin/zara-jira-mcp

VOLUME /data
ENV PM_MEMORY_DB_PATH=/data/pm_memory.db

ENTRYPOINT ["zara-jira-mcp"]
