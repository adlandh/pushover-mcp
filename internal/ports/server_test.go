package ports

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/adlandh/pushover-mcp/internal/application"
	"github.com/adlandh/pushover-mcp/internal/domain"
)

type fakeNotificationSender struct {
	err          error
	called       bool
	notification domain.Notification
}

func (f *fakeNotificationSender) Send(_ context.Context, notification domain.Notification) error {
	f.called = true
	f.notification = notification
	return f.err
}

func TestNewServer_RegistersSendTool(t *testing.T) {
	useCase := application.NewSendNotificationUseCase(&fakeNotificationSender{})
	s := NewServer("test-server", "1.0.0", useCase)

	tool := s.GetTool("send")
	if tool == nil {
		t.Fatal("send tool was not registered")
	}

	if tool.Tool.Name != "send" {
		t.Fatalf("tool name = %q, want %q", tool.Tool.Name, "send")
	}
	if tool.Tool.Description != "Sends a notification via Pushover." {
		t.Fatalf("tool description = %q", tool.Tool.Description)
	}

	if tool.Tool.InputSchema.Type != "object" {
		t.Fatalf("input schema type = %q, want %q", tool.Tool.InputSchema.Type, "object")
	}

	if len(tool.Tool.InputSchema.Required) != 1 || tool.Tool.InputSchema.Required[0] != "message" {
		t.Fatalf("required fields = %v, want [message]", tool.Tool.InputSchema.Required)
	}

	additional, ok := tool.Tool.InputSchema.AdditionalProperties.(bool)
	if !ok || additional {
		t.Fatalf("additionalProperties = %v, want false", tool.Tool.InputSchema.AdditionalProperties)
	}
}

func TestSendToolHandler_Success(t *testing.T) {
	fakeSender := &fakeNotificationSender{}
	useCase := application.NewSendNotificationUseCase(fakeSender)
	s := NewServer("test-server", "1.0.0", useCase)

	tool := s.GetTool("send")
	if tool == nil {
		t.Fatal("send tool was not registered")
	}

	priority := 1
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "send",
			Arguments: map[string]any{
				"message":  "hello",
				"priority": priority,
				"url":      "https://example.com",
			},
		},
	}

	result, err := tool.Handler(context.Background(), request)
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.IsError {
		t.Fatalf("result IsError = true, want false")
	}

	if !fakeSender.called {
		t.Fatal("sender.Send was not called")
	}
	if fakeSender.notification.Message != "hello" {
		t.Fatalf("message = %q, want %q", fakeSender.notification.Message, "hello")
	}
	if fakeSender.notification.Priority == nil || *fakeSender.notification.Priority != 1 {
		t.Fatalf("priority = %v, want 1", fakeSender.notification.Priority)
	}
	if fakeSender.notification.URL == nil || *fakeSender.notification.URL != "https://example.com" {
		t.Fatalf("url = %v, want https://example.com", fakeSender.notification.URL)
	}

	if len(result.Content) == 0 {
		t.Fatal("result content is empty")
	}
	text := mcp.GetTextFromContent(result.Content[0])
	if text != "Notification sent." {
		t.Fatalf("result text = %q, want %q", text, "Notification sent.")
	}
}

func TestSendToolHandler_InvalidArguments(t *testing.T) {
	useCase := application.NewSendNotificationUseCase(&fakeNotificationSender{})
	s := NewServer("test-server", "1.0.0", useCase)

	tool := s.GetTool("send")
	if tool == nil {
		t.Fatal("send tool was not registered")
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "send",
			Arguments: map[string]any{
				"message":  "hello",
				"priority": "high",
			},
		},
	}

	result, err := tool.Handler(context.Background(), request)
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if !result.IsError {
		t.Fatal("result IsError = false, want true")
	}

	if len(result.Content) == 0 {
		t.Fatal("result content is empty")
	}
	text := mcp.GetTextFromContent(result.Content[0])
	if !strings.Contains(text, "invalid tool arguments") {
		t.Fatalf("result text = %q, want to contain %q", text, "invalid tool arguments")
	}
}

func TestSendToolHandler_UseCaseError(t *testing.T) {
	fakeSender := &fakeNotificationSender{err: errors.New("pushover unavailable")}
	useCase := application.NewSendNotificationUseCase(fakeSender)
	s := NewServer("test-server", "1.0.0", useCase)

	tool := s.GetTool("send")
	if tool == nil {
		t.Fatal("send tool was not registered")
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "send",
			Arguments: map[string]any{
				"message": "hello",
			},
		},
	}

	result, err := tool.Handler(context.Background(), request)
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if !result.IsError {
		t.Fatal("result IsError = false, want true")
	}

	if len(result.Content) == 0 {
		t.Fatal("result content is empty")
	}
	text := mcp.GetTextFromContent(result.Content[0])
	if !strings.Contains(text, "Failed to send notification") {
		t.Fatalf("result text = %q, want to contain %q", text, "Failed to send notification")
	}
}
