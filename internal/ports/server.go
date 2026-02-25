package ports

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/domain"
)

type sendArguments struct {
	Title    *string `json:"title,omitempty"`
	Priority *int    `json:"priority,omitempty"`
	Sound    *string `json:"sound,omitempty"`
	URL      *string `json:"url,omitempty"`
	URLTitle *string `json:"url_title,omitempty"`
	Device   *string `json:"device,omitempty"`
	Message  string  `json:"message"`
}

func NewServer(name, version string, useCase *application.SendNotificationUseCase) *server.MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s.AddTool(buildSendTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args sendArguments
		if err := request.BindArguments(&args); err != nil {
			return mcp.NewToolResultErrorf("invalid tool arguments: %v", err), nil
		}

		notification := domain.Notification{
			Message:  args.Message,
			Title:    args.Title,
			Priority: args.Priority,
			Sound:    args.Sound,
			URL:      args.URL,
			URLTitle: args.URLTitle,
			Device:   args.Device,
		}

		if err := useCase.Execute(ctx, notification); err != nil {
			return mcp.NewToolResultErrorf("Failed to send notification: %v", err), nil
		}

		return mcp.NewToolResultText("Notification sent."), nil
	})

	return s
}

func buildSendTool() mcp.Tool {
	return mcp.NewTool("send",
		mcp.WithDescription("Sends a notification via Pushover."),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to send"),
		),
		mcp.WithString("title",
			mcp.Description("Message title"),
		),
		mcp.WithNumber("priority",
			mcp.Description("Priority from -2 to 2 (-2: lowest, 2: emergency)"),
			mcp.Min(-2),
			mcp.Max(2),
		),
		mcp.WithString("sound",
			mcp.Description("Notification sound"),
		),
		mcp.WithString("url",
			mcp.Description("URL to include"),
		),
		mcp.WithString("url_title",
			mcp.Description("Title for the URL"),
		),
		mcp.WithString("device",
			mcp.Description("Target specific device"),
		),
		mcp.WithSchemaAdditionalProperties(false),
	)
}
