package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type SavedReceiverService interface {
	AddInternalReceiver(ctx *gin.Context, receiver model.InternalReceiver) error
	AddExternalReceiver(ctx *gin.Context, receiver model.ExternalReceiver) error
}
