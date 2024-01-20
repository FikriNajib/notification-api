package notification

import (
	"context"
	"notification-api/domain"
)

type EmailRepository interface {
	PostEmail(ctx context.Context, request domain.MetadataEmail) (domain.ResponseNotif, error)
}
