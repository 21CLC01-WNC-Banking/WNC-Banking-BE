package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type SavedReceiverService interface {
	AddReceiver(ctx *gin.Context, receiver model.Receiver) error
	GetAllReceivers(ctx *gin.Context) (*[]model.SavedReceiverResponse, error)
	UpdateNickname(ctx *gin.Context, id int64, newNickname string) error
	DeleteReceiver(ctx *gin.Context, id int64) error
}
