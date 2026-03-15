# AGENTS.md

## Scope

These instructions apply to the entire `pushover-mcp` repository.

The project is a small Go service that exposes Pushover notification sending through an MCP server. Keep changes minimal, explicit, and aligned with the current hexagonal structure.

## Repository Snapshot

- Module: `github.com/adlandh/pushover-mcp`
- Go version: `1.26`
- Entry point: `main.go`
- Core dependency: `github.com/mark3labs/mcp-go`
- Configuration library: `github.com/caarlos0/env/v11`

## Architecture

Preserve the existing dependency direction:

- `internal/domain`: core models and interfaces
- `internal/application`: use cases and orchestration
- `internal/adapters`: external integrations such as the Pushover HTTP client
- `internal/config`: environment parsing and config loading
- `internal/ports`: MCP server wiring and transport-facing handlers

Rule:
- Dependencies should point inward. `domain` must not depend on application, adapters, config, or ports.
- Keep MCP-specific concerns in `ports`.
- Keep HTTP and Pushover API details in `adapters`.
- Put business flow in `application`, not in handlers or clients.

## Working Style

- Prefer small, targeted edits over broad refactors.
- Match the surrounding code before introducing new patterns.
- Do not move code across layers without a clear architectural reason.
- When behavior changes, update or add tests in the same pass.
- Avoid speculative abstractions unless duplication is already real.

## Build, Test, and Verification

Run the narrowest command that proves the change first, then broader validation if needed.

```bash
# Run all tests
go test ./...

# Run a single package test
go test ./internal/adapters/... -run TestNewClient_Validation -v

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Format Go code
gofmt -w .

# Build binary
go build -o pushover-mcp .

# Install locally
go install .

# Run linter when available
golangci-lint run ./...
```

Expectations:

- At minimum, run tests for the package you changed.
- Run `go test ./...` when changes cross package boundaries or affect wiring.
- Run `gofmt -w` on every modified Go file.
- Run `golangci-lint run ./...` before finishing if the environment has it installed.

## Code Style

### Imports

Use import groups separated by blank lines:

1. Standard library
2. Third-party packages
3. Internal project imports

Example:

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/domain"
)
```

### Formatting

- Use `gofmt` for formatting.
- Use tabs for indentation.
- Keep lines under 120 characters when practical.
- Prefer short functions with early returns.

### Naming

- Packages: lowercase, concise, single-purpose
- Exported identifiers: PascalCase
- Unexported identifiers: camelCase
- Constants: descriptive names, grouped in `const` blocks when related
- Interfaces: use descriptive nouns, often ending in `-er`
- Test functions: `Test<FunctionName>_<Scenario>`
- Test doubles/helpers: prefix with `fake`

### Error Handling

Wrap returned errors with context using `fmt.Errorf` and `%w`.

```go
if err != nil {
	return fmt.Errorf("create request: %w", err)
}
```

Validation errors should be direct and specific.

```go
if cfg.APIToken == "" {
	return nil, fmt.Errorf("missing APIToken")
}
```

Rules:

- Do not discard underlying errors when they add debugging value.
- Use lowercase error messages without trailing punctuation.
- Add context that identifies the failed operation.

### Struct Layout

Order fields to reduce padding where it does not harm readability.

```go
type PushoverClient struct {
	httpClient *http.Client
	apiToken   string
	userKey    string
	apiURL     string
}
```

### Control Flow and Whitespace

The repo uses stricter whitespace and readability conventions:

- Add a blank line before `if`, `for`, `switch`, and `return` when it improves readability.
- Add a blank line after a function declaration before the first statement.
- Keep tightly related statements together; do not add whitespace mechanically.

### Complexity

- Target cyclomatic complexity under `10`.
- Split branches into helpers before nesting becomes hard to scan.
- Prefer explicit logic over clever compression.

## Testing Guidelines

- Prefer table-driven tests.
- Use `t.Run` for named scenarios.
- Mock external HTTP behavior instead of hitting real services.
- Keep tests deterministic and isolated.
- In config tests, set environment variables explicitly.
- Add regression tests for bug fixes.

Example shape:

```go
func TestNewClient_Validation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		// ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ...
		})
	}
}
```

## Linting Expectations

The repo expects compatibility with these linters:

- `wrapcheck`
- `wsl_v5`
- `cyclop`
- `gosec`
- `errcheck`
- `gosmopolitan`

Write code that satisfies them by default, especially around error wrapping, unchecked errors, whitespace, and function complexity.

## File-Level Guidance

### `internal/domain`

- Keep it small and dependency-free.
- Only domain concepts, value objects, and interfaces belong here.

### `internal/application`

- Express use cases in terms of domain interfaces.
- Avoid direct knowledge of HTTP, env parsing, or MCP transport details.

### `internal/adapters`

- Encapsulate Pushover API request construction, response handling, and HTTP client behavior.
- Keep serialization and protocol details here.

### `internal/config`

- Handle environment-driven configuration only.
- Validate required settings close to loading.

### `internal/ports`

- Define MCP tools, handlers, and server assembly.
- Translate MCP input/output into application calls without embedding business logic.

## Change Checklist

Before finishing:

- Code is formatted with `gofmt`
- Imports are grouped correctly
- Errors are wrapped with useful context
- New behavior has tests, or existing tests were updated
- Layer boundaries still hold
- The narrowest relevant tests were run
- `golangci-lint run ./...` was run if available
