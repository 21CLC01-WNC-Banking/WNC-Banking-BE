package bean

import (
	"context"
	"time"
)

type MailClient interface {
	GenerateOTPBody(to, code, context string, ttl time.Duration) string
	GenerateRandomPasswordBody(to, password string) string
	SendEmail(ctx context.Context, to, subject, body string) error
}
