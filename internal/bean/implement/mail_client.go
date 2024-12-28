package beanimplement

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"gopkg.in/gomail.v2"
)

type MailClient struct {
	from   string
	dialer *gomail.Dialer
}

func NewMailClient() bean.MailClient {
	host, _ := env.GetEnv("MAIL_HOST")
	portStr, _ := env.GetEnv("MAIL_PORT")
	port, _ := strconv.Atoi(portStr)
	username, _ := env.GetEnv("MAIL_USERNAME")
	password, _ := env.GetEnv("MAIL_PASSWORD")
	from, _ := env.GetEnv("MAIL_FROM")

	dialer := gomail.NewDialer(host, port, username, password)

	return &MailClient{
		from:   from,
		dialer: dialer,
	}
}

func (m *MailClient) GenerateOTPBody(to, code, context string, ttl time.Duration) string {
	contextMessage := fmt.Sprintf("Nhập lại mã sau để %s:", context)

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Bạn đã chọn <strong>%s</strong> làm địa chỉ email của mình.</p>
			<p>%s</p>
			<h2 style="color: #4CAF50;">%s</h2>
			<p>Mã này sẽ hết hạn sau %d phút kể từ khi email này được gửi.</p>
			<hr>
			<p style="color: red;"><strong>Cảnh báo: Không bao giờ chia sẻ mã OTP của bạn với bất kỳ ai.</strong></p>
			<p><strong>Tại sao bạn nhận được email này?</strong></p>
			<p>Chúng tôi yêu cầu xác minh bất cứ khi nào mật khẩu được cập nhật hoặc giao dịch được thực hiện.</p>
			<p>Nếu bạn không thực hiện yêu cầu này, bạn có thể bỏ qua email này.</p>
		</body>
		</html>`, to, contextMessage, code, ttl/time.Minute)
	return body
}

func (m *MailClient) GenerateRandomPasswordBody(to, password string) string {
	contextMessage := "Mật khẩu để đăng nhập vào tài khoản của bạn"

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="vi">
		<body>
			<p>Bạn đã chọn <strong>%s</strong> làm địa chỉ email của mình.</p>
			<p>%s</p>
			<p>Mật khẩu tạm thời của bạn là:</p>
			<h2 style="color: #4CAF50;">%s</h2>
			<hr>
			<p style="color: red;"><strong>Cảnh báo: Không bao giờ chia sẻ mật khẩu của bạn với bất kỳ ai.</strong></p>
			<p><strong>Tại sao bạn nhận được email này?</strong></p>
			<p>Chúng tôi yêu cầu xác minh bất cứ khi nào mật khẩu được cập nhật hoặc giao dịch được thực hiện.</p>
			<p>Nếu bạn không thực hiện yêu cầu này, bạn có thể bỏ qua email này.</p>
		</body>
		</html>
	`, to, contextMessage, password)
	return body
}

func (m *MailClient) SendEmail(ctx context.Context, to, subject, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", m.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return m.dialer.DialAndSend(message)
}
