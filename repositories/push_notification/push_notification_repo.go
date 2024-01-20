package push_notification

import (
	"context"
	"notification-api/domain"
)

type PushNotifRepository interface {
	PostPushNotif(ctx context.Context, request domain.PushNotificationRequest) (domain.ResponseNotif, error)
}
