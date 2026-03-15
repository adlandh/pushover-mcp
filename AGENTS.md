# AGENTS.md

## Scope

These instructions apply to the entire `pushover-mcp` repository.

`pushover-mcp` is a small Go service that exposes one MCP tool for sending Pushover notifications. Keep changes minimal and aligned with the current hexagonal structure.

## Repository Snapshot

- Module: `github.com/adlandh/pushover-mcp`
- Go version: `1.26`
- Entry point: `main.go`
- MCP library: `github.com/mark3labs/mcp-go`
- Config library: `github.com/caarlos0/env/v11`

## Architecture

Preserve dependency direction:

- `internal/domain`: core models and interfaces
- `internal/application`: use cases and orchestration
- `internal/adapters`: external integrations such as the Pushover HTTP client
- `internal/config`: environment parsing and config loading
- `internal/ports`: MCP server wiring and transport-facing handlers

Rules:

- Dependencies point inward.
- `domain` must not depend on `application`, `adapters`, `config`, or `ports`.
- Keep HTTP and Pushover API details in `adapters`.
- Keep MCP-specific concerns in `ports`.
- Keep validation and orchestration in `application`.
- Keep `main.go` thin.

## Working Style

- Prefer small, targeted edits over broad refactors.
- Match existing patterns before introducing new ones.
- Do not move code across layers without a clear reason.
- When behavior changes, update or add tests in the same pass.
- Avoid speculative abstractions.

## Build, Test, and Verification

Use the narrowest command that proves the change first.

```bash
# Run all tests
go test ./...

# Run a single test
go test ./internal/adapters/... -run TestNewClient_Validation -v

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Format Go code
goimports -w .

# Build binary
go build -o pushover-mcp .

# Install locally
go install .

# Run linter
golangci-lint run ./...
```

Expectations:

- Run tests for the package you changed at minimum.
- Run `go test ./...` when changes cross package boundaries or affect wiring.
- Run `goimports -w` on every modified Go file.
- Run `golangci-lint run ./...` before finishing if available.

## Code Style

### Imports

Use import groups separated by blank lines:
1. Standard library
2. Third-party packages
3. Internal project imports

### Formatting

- Use `goimports` for formatting and import normalization.
- Use tabs for indentation.
- Keep lines under 120 characters when practical.
- Prefer short functions with early returns.

### Naming

- Packages: lowercase, concise, single-purpose
- Exported identifiers: PascalCase
- Unexported identifiers: camelCase
- Constants: descriptive names in `const` blocks when related
- Interfaces: descriptive nouns, often ending in `-er`
- Test functions: `Test<FunctionName>_<Scenario>`
- Test doubles/helpers: prefix with `fake`

### Error Handling

Wrap returned errors with context using `fmt.Errorf` and `%w`.

```go
if err != nil {
	return fmt.Errorf("create request: %w", err)
}
```

Rules:

- Validation errors should be direct and specific.
- Use lowercase error messages without trailing punctuation.
- Do not discard underlying errors when they add debugging value.
- Add operation context to wrapped errors.

### Struct Layout and Whitespace

- Order struct fields to reduce padding when it does not harm readability.
- `wsl_v5` is enabled; add blank lines where they improve readability.
- Add a blank line after a function declaration before the first statement.
- Follow surrounding style and rerun the linter.

### Complexity

- Target cyclomatic complexity under `10`.
- Split complex logic into helpers before nesting becomes hard to scan.
- Prefer explicit logic over clever compression.

## Testing Guidelines

- Prefer table-driven tests.
- Use `t.Run` for named scenarios.
- Keep tests deterministic and isolated.
- Mock external HTTP behavior instead of hitting real services.
- In config tests, set environment variables explicitly.
- Add regression tests for bug fixes.
- Extract repeated test literals into constants when lint or review points them out.

Useful repo patterns:

- Adapter tests use `httptest.Server`.
- Application tests use `fake` senders.
- Port tests can call MCP handlers directly.

## Linting Expectations

Write code compatible with these linters:

- `wrapcheck`
- `wsl_v5`
- `cyclop`
- `gosec`
- `errcheck`
- `gosmopolitan`
- `govet`

## File-Level Guidance

- `internal/domain`: keep it dependency-free.
- `internal/application`: express use cases in terms of domain interfaces.
- `internal/adapters`: encapsulate Pushover HTTP behavior.
- `internal/config`: handle env-driven configuration only.
- `internal/ports`: translate MCP input/output without embedding business logic.

## Agent Notes

- No Cursor rules or Copilot instruction files are present in this repository.
- Do not introduce new tooling or config formats unless clearly needed.
- Keep `README.md` examples aligned with actual runtime behavior and config keys.

## Change Checklist

Before finishing:

- Code is formatted with `goimports`
- Imports are grouped correctly
- Errors are wrapped with useful context
- New behavior has tests, or existing tests were updated
- Layer boundaries still hold
- The narrowest relevant tests were run
- `golangci-lint run ./...` was run if available
