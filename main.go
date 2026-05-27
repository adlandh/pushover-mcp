package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/config"
	"github.com/adlandh/pushover-mcp/internal/driven"
	"github.com/adlandh/pushover-mcp/internal/driver"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "pushover-mcp"
	serverVersion = "1.0.0"
)

func buildServer(env config.EnvConfig) (*server.MCPServer, error) {
	httpClient := &http.Client{Timeout: env.Timeout}

	sender, err := driven.NewPushoverClient(env.Pushover, httpClient)
	if err != nil {
		return nil, fmt.Errorf("error creating sender: %w", err)
	}

	useCase := application.NewSendNotificationUseCase(sender)

	return driver.NewServer(serverName, serverVersion, useCase), nil
}

func run() error {
	env, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	mcpServer, err := buildServer(env)
	if err != nil {
		return err
	}

	if err := server.ServeStdio(mcpServer); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
