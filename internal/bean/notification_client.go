package bean

import "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"

type NotificationClient interface {
	Send()
	SendTest(req model.NotificationRequest)
}