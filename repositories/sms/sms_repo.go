package sms

import (
	"context"
	"notification-api/domain"
)

type SmsRepository interface {
	SendSMS(ctx context.Context, recipient, otp, otpType string) (domain.ResponseNotif, error)
}
