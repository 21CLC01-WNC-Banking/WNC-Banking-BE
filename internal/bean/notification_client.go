package bean

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type NotificationClient interface {
	SendTest(req model.NotificationRequest)
	SendAndSave(ctx *gin.Context, req interface{})
}
