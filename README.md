# Pushover MCP

[![Go Reference](https://pkg.go.dev/badge/github.com/adlandh/context-logger.svg)](https://pkg.go.dev/github.com/adlandh/context-logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/adlandh/context-logger)](https://goreportcard.com/report/github.com/adlandh/context-logger)

MCP service with one tool: `send`.

The service sends notifications through [Pushover](https://pushover.net/).

## Requirements

- Go 1.26+
- Pushover app token
- Pushover user key

## Environment variables

- `PUSHOVER_API_TOKEN` - required
- `PUSHOVER_USER_KEY` - required
- `PUSHOVER_API_URL` - optional (default: `https://api.pushover.net/1/messages.json`)
- `PUSHOVER_TIMEOUT` - optional HTTP timeout as Go duration (default: `15s`, examples: `5s`, `30s`, `1m`)

## Install

```bash
go install github.com/adlandh/pushover-mcp@latest
```

## MCP client setup

Install first:

```bash
go install github.com/adlandh/pushover-mcp@latest
```

The binary will be available as `pushover-mcp` in your `GOBIN` (or `$(go env GOPATH)/bin` if `GOBIN` is not set).
Use absolute path to that binary in client configs.

### Codex

Add server config to your Codex MCP settings in TOML format:

```toml
[mcp_servers.pushover]
command = "/absolute/path/to/pushover-mcp"

[mcp_servers.pushover.env]
PUSHOVER_API_TOKEN = "YOUR_TOKEN"
PUSHOVER_USER_KEY = "YOUR_USER_KEY"
```

### Claude Desktop

Edit `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "pushover": {
      "command": "/absolute/path/to/pushover-mcp",
      "env": {
        "PUSHOVER_API_TOKEN": "YOUR_TOKEN",
        "PUSHOVER_USER_KEY": "YOUR_USER_KEY"
      }
    }
  }
}
```

## Request structure

Tool name: `send`

Arguments payload:

```json
{
  "message": "Deploy finished", 
  "title": "CI",
  "priority": 0,
  "sound": "pushover",
  "url": "https://example.com/build/123",
  "url_title": "Open build",
  "device": "iphone"
}
```

MCP `tools/call` request example:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "send",
    "arguments": {
      "message": "Deploy finished",
      "title": "CI",
      "priority": 0
    }
  }
}
```

Typical paths:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%/Claude/claude_desktop_config.json`

### OpenCode

Add the same MCP server block to your OpenCode config:

```json
{
  "mcpServers": {
    "pushover": {
      "command": "/absolute/path/to/pushover-mcp",
      "env": {
        "PUSHOVER_API_TOKEN": "YOUR_TOKEN",
        "PUSHOVER_USER_KEY": "YOUR_USER_KEY"
      }
    }
  }
}
```

If your OpenCode setup uses per-project config, put this block into the project-level MCP config instead of global config.

## Quick local check (bash)

You can ping the tool directly from bash:

```bash
{
  printf '%s\n' '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"bash-test","version":"0.1.0"}}}'
  printf '%s\n' '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'
  printf '%s\n' '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"send","arguments":{"message":"Test from bash","title":"MCP test","priority":0}}}'
} | ./pushover-mcp
```
