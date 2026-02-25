package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

type SendNotificationUseCase struct {
	sender domain.NotificationSender
}

func NewSendNotificationUseCase(sender domain.NotificationSender) *SendNotificationUseCase {
	return &SendNotificationUseCase{sender: sender}
}

func (u *SendNotificationUseCase) Execute(ctx context.Context, notification domain.Notification) error {
	if strings.TrimSpace(notification.Message) == "" {
		return errors.New("message is required")
	}

	if notification.Priority != nil {
		if *notification.Priority < -2 || *notification.Priority > 2 {
			return errors.New("priority must be between -2 and 2")
		}
	}

	if err := u.sender.Send(ctx, notification); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	return nil
}
