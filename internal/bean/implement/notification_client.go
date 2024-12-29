package beanimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/websocket"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
)

type NotificationClient struct {
	client *websocket.Server
}

func NewNotificationClient(client *websocket.Server) bean.NotificationClient {

	return &NotificationClient{
		client: client,
	}
}

func (c *NotificationClient) SendTest(req model.NotificationRequest) {
	c.client.SendToDevice(req.DeviceId, req.Title+"\n"+req.Content)
}
