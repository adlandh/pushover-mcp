package main

import (
	"log"
	"net/http"

	"github.com/adlandh/pushover-mcp/internal/adapters"
	"github.com/adlandh/pushover-mcp/internal/ports"
	"github.com/mark3labs/mcp-go/server"

	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/config"
)

const (
	serverName    = "pushover-mcp"
	serverVersion = "1.0.0"
)

func main() {
	env, err := config.FromEnv()
	if err != nil {
		log.Fatalf("configuration error: %v\n", err)
	}

	httpClient := &http.Client{Timeout: env.Timeout}
	if sender, err := adapters.NewPushoverClient(env.Pushover, httpClient); err != nil {
		log.Fatalf("error creating sender: %v\n", err)
	} else {
		useCase := application.NewSendNotificationUseCase(sender)
		mcpServer := ports.NewServer(serverName, serverVersion, useCase)

		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("error starting server: %v\n", err)
		}
	}
}
