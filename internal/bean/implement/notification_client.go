package beanimplement

import (
	"encoding/json"
	"fmt"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/websocket"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/noti"
	"github.com/gin-gonic/gin"
)

type NotificationClient struct {
	client                 *websocket.Server
	notificationRepository repository.NotificationRepository
}

func NewNotificationClient(client *websocket.Server, notificationRepository repository.NotificationRepository) bean.NotificationClient {

	return &NotificationClient{
		client:                 client,
		notificationRepository: notificationRepository,
	}
}

func (c *NotificationClient) SendTest(req model.NotificationRequest) {
	c.client.SendToDevice(req.DeviceId, req.Title+"\n"+req.Content)
}

func (c *NotificationClient) SendAndSave(ctx *gin.Context, req model.NotificationResponse) {
	// Save notification
	content, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	notification := &entity.Notification{
		Type:    req.Type,
		Content: string(content),
		IsSeen:  false,
		UserID:  int64(req.DeviceId),
	}
	err = c.notificationRepository.CreateCommand(ctx, notification)
	if err != nil {
		fmt.Println(err)
	}

	// Send notification
	c.client.SendToDevice(req.DeviceId, noti.GenerateContent(req))
}
