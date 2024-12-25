package bean

import (
	"context"
	"time"
)

type MailClient interface {
	GenerateOTPBody(to, code, context string, ttl time.Duration) string
	SendEmail(ctx context.Context, to, subject, body string) error
}
