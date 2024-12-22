package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type TransactionService interface {
	PreInternalTransfer(ctx *gin.Context, transferReq model.PreInternalTransferRequest) (string, error)
	SendOTPToEmail(ctx *gin.Context, email string, transactionId string) error
	InternalTransfer(ctx *gin.Context, transferReq model.InternalTransferRequest) (*entity.Transaction, error)
}
