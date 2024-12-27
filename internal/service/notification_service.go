package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type NotificationService interface {
	GetAllNotifications(ctx *gin.Context, userId int64) ([]entity.Notification, error)
	PatchNotification(ctx *gin.Context, userId int64, id int64) error
}
