package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

type fakeSender struct {
	err          error
	called       bool
	notification domain.Notification
}

func (f *fakeSender) Send(_ context.Context, notification domain.Notification) error {
	f.called = true
	f.notification = notification
	return f.err
}

func TestSendNotificationUseCase_Execute_Success(t *testing.T) {
	sender := &fakeSender{}
	useCase := NewSendNotificationUseCase(sender)

	priority := 1
	title := "Test"
	notification := domain.Notification{
		Message:  "hello",
		Title:    &title,
		Priority: &priority,
	}

	err := useCase.Execute(context.Background(), notification)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !sender.called {
		t.Fatal("sender.Send was not called")
	}
	if sender.notification.Message != "hello" {
		t.Fatalf("message = %q, want %q", sender.notification.Message, "hello")
	}
	if sender.notification.Title == nil || *sender.notification.Title != "Test" {
		t.Fatalf("title = %v, want Test", sender.notification.Title)
	}
	if sender.notification.Priority == nil || *sender.notification.Priority != 1 {
		t.Fatalf("priority = %v, want 1", sender.notification.Priority)
	}
}

func TestSendNotificationUseCase_Execute_AllOptionalFields(t *testing.T) {
	sender := &fakeSender{}
	useCase := NewSendNotificationUseCase(sender)

	priority := 0
	title := "Title"
	sound := "pushover"
	url := "https://example.com"
	urlTitle := "Link"
	device := "iphone"

	notification := domain.Notification{
		Message:  "hello",
		Title:    &title,
		Priority: &priority,
		Sound:    &sound,
		URL:      &url,
		URLTitle: &urlTitle,
		Device:   &device,
	}

	err := useCase.Execute(context.Background(), notification)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !sender.called {
		t.Fatal("sender.Send was not called")
	}
	if sender.notification.Sound == nil || *sender.notification.Sound != "pushover" {
		t.Fatalf("sound = %v, want pushover", sender.notification.Sound)
	}
	if sender.notification.URL == nil || *sender.notification.URL != "https://example.com" {
		t.Fatalf("url = %v, want https://example.com", sender.notification.URL)
	}
	if sender.notification.Device == nil || *sender.notification.Device != "iphone" {
		t.Fatalf("device = %v, want iphone", sender.notification.Device)
	}
}

func TestSendNotificationUseCase_Execute_MessageRequired(t *testing.T) {
	sender := &fakeSender{}
	useCase := NewSendNotificationUseCase(sender)

	err := useCase.Execute(context.Background(), domain.Notification{Message: "   "})
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if err.Error() != "message is required" {
		t.Fatalf("error = %q, want %q", err.Error(), "message is required")
	}
	if sender.called {
		t.Fatal("sender.Send was called, want not called")
	}
}

func TestSendNotificationUseCase_Execute_PriorityOutOfRange(t *testing.T) {
	tests := []struct {
		name     string
		priority int
	}{
		{name: "below min", priority: -3},
		{name: "above max", priority: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sender := &fakeSender{}
			useCase := NewSendNotificationUseCase(sender)

			err := useCase.Execute(context.Background(), domain.Notification{
				Message:  "hello",
				Priority: &tc.priority,
			})
			if err == nil {
				t.Fatal("Execute() error = nil, want non-nil")
			}
			if err.Error() != "priority must be between -2 and 2" {
				t.Fatalf("error = %q, want %q", err.Error(), "priority must be between -2 and 2")
			}
			if sender.called {
				t.Fatal("sender.Send was called, want not called")
			}
		})
	}
}

func TestSendNotificationUseCase_Execute_SenderError(t *testing.T) {
	sender := &fakeSender{err: errors.New("network error")}
	useCase := NewSendNotificationUseCase(sender)

	err := useCase.Execute(context.Background(), domain.Notification{Message: "hello"})
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("error = %q, want to contain %q", err.Error(), "network error")
	}
	if !sender.called {
		t.Fatal("sender.Send was not called")
	}
}
