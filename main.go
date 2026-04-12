package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adlandh/pushover-mcp/internal/adapters"
	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/config"
	"github.com/adlandh/pushover-mcp/internal/ports"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "pushover-mcp"
	serverVersion = "1.0.0"
)

func run() error {
	env, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	httpClient := &http.Client{Timeout: env.Timeout}

	sender, err := adapters.NewPushoverClient(env.Pushover, httpClient)
	if err != nil {
		return fmt.Errorf("error creating sender: %w", err)
	}

	useCase := application.NewSendNotificationUseCase(sender)
	mcpServer := ports.NewServer(serverName, serverVersion, useCase)

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
