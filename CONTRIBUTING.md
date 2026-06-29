# Contributing

Thanks for your interest in contributing to zara-jira-mcp.

## Development Setup

```bash
git clone https://github.com/aldok10/zara-jira-mcp.git
cd zara-jira-mcp
cp .env.example .env  # edit with your credentials
make build
make test
```

## Architecture

```
cmd/server/          Entry point (uber-go/fx DI)
config/              Env-based configuration
domain/              Interfaces + models (never import internal/)
internal/            Implementations
application/tools/   MCP tool handlers
transport/           Tool registration + MCP server
```

## Adding a New Tool

1. Add handler method in `application/tools/`
2. Register in the appropriate `transport/*.go` file
3. Add to the correct module in `transport/server.go` `modules` map
4. Update `SKILL.md` with params and usage

## Code Style

- Go standard formatting (`gofmt`)
- Uber Go style guide conventions
- Keep handlers focused — one tool, one function
- Error messages should be user-friendly (PM/SM reads them)
- No panics in handlers, always return error results

## Pull Requests

- One feature per PR
- Include test if adding logic
- Update SKILL.md if changing tool signatures
- Keep commits focused and messages clear

## Issues

- Bug reports: include PM_PROFILE, client name, and error output
- Feature requests: describe the PM/SM workflow it enables

## License

By contributing, you agree your code is licensed under MIT.
