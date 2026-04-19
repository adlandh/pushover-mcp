package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

const (
	testMessage         = "hello"
	errMessageRequired  = "message is required"
	errPriorityOutRange = "priority must be between -2 and 2"
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

func assertStringPtr(t *testing.T, got *string, want, field string) {
	t.Helper()

	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %q", field, got, want)
	}
}

func assertIntPtr(t *testing.T, got *int, want int, field string) {
	t.Helper()

	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %d", field, got, want)
	}
}

func assertValidationError(t *testing.T, sender *fakeSender, err error, want string) {
	t.Helper()

	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}

	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}

	if sender.called {
		t.Fatal("sender.Send was called, want not called")
	}
}

func TestSendNotificationUseCase_Execute_Success(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	priority := 1
	title := "Test"
	notification := domain.Notification{
		Message:  testMessage,
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

	if sender.notification.Message != testMessage {
		t.Fatalf("message = %q, want %q", sender.notification.Message, testMessage)
	}

	assertStringPtr(t, sender.notification.Title, "Test", "title")
	assertIntPtr(t, sender.notification.Priority, 1, "priority")
}

func TestSendNotificationUseCase_Execute_AllOptionalFields(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	priority := 0
	title := "Title"
	sound := "pushover"
	url := "https://example.com"
	urlTitle := "Link"
	device := "iphone"

	notification := domain.Notification{
		Message:  testMessage,
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

	assertStringPtr(t, sender.notification.Sound, sound, "sound")
	assertStringPtr(t, sender.notification.URL, url, "url")
	assertStringPtr(t, sender.notification.Device, device, "device")
}

func TestSendNotificationUseCase_Execute_MessageRequired(t *testing.T) {
	sender, useCase := newUseCaseWithFake()

	err := useCase.Execute(context.Background(), domain.Notification{Message: "   "})
	assertValidationError(t, sender, err, errMessageRequired)
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
			assertValidationError(t, sender, err, errPriorityOutRange)
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
