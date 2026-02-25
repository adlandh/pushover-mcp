package domain

import (
	"context"
)

type NotificationSender interface {
	Send(ctx context.Context, notification Notification) error
}
