package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

const (
	testMessage = "hello"
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

func newUseCaseWithFake() (*fakeSender, *SendNotificationUseCase) {
	sender := &fakeSender{}
	useCase := NewSendNotificationUseCase(sender)

	return sender, useCase
}

func assertString(t *testing.T, got, want, field string) {
	t.Helper()

	if got != want {
		t.Fatalf("%s = %q, want %q", field, got, want)
	}
}

func assertIntPtr(t *testing.T, got *int, want int, field string) {
	t.Helper()

	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %d", field, got, want)
	}
}

func assertValidationError(t *testing.T, sender *fakeSender, err, want error) {
	t.Helper()

	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}

	if sender.called {
		t.Fatal("sender.Send was called, want not called")
	}
}

func TestSendNotificationUseCase_Execute_Success(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	priority := 1
	notification := domain.Notification{
		Message:  testMessage,
		Title:    "Test",
		Priority: &priority,
	}

	err := useCase.Execute(context.Background(), notification)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !sender.called {
		t.Fatal("sender.Send was not called")
	}

	if sender.notification.Message != testMessage {
		t.Fatalf("message = %q, want %q", sender.notification.Message, testMessage)
	}

	assertString(t, sender.notification.Title, "Test", "title")
	assertIntPtr(t, sender.notification.Priority, 1, "priority")
}

func TestSendNotificationUseCase_Execute_AllOptionalFields(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	priority := 0
	sound := "pushover"
	url := "https://example.com"
	device := "iphone"

	notification := domain.Notification{
		Message:  testMessage,
		Title:    "Title",
		Priority: &priority,
		Sound:    sound,
		URL:      url,
		URLTitle: "Link",
		Device:   device,
	}

	err := useCase.Execute(context.Background(), notification)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !sender.called {
		t.Fatal("sender.Send was not called")
	}

	assertString(t, sender.notification.Sound, sound, "sound")
	assertString(t, sender.notification.URL, url, "url")
	assertString(t, sender.notification.Device, device, "device")
}

func TestSendNotificationUseCase_Execute_MessageRequired(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	err := useCase.Execute(context.Background(), domain.Notification{Message: "   "})
	assertValidationError(t, sender, err, ErrMessageRequired)
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
			sender, useCase := newUseCaseWithFake()

			err := useCase.Execute(context.Background(), domain.Notification{
				Message:  testMessage,
				Priority: &tc.priority,
			})
			assertValidationError(t, sender, err, ErrPriorityOutRange)
		})
	}
}

func TestSendNotificationUseCase_Execute_SenderError(t *testing.T) {
	sender := &fakeSender{err: errors.New("network error")}
	useCase := NewSendNotificationUseCase(sender)

	err := useCase.Execute(context.Background(), domain.Notification{Message: testMessage})
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
