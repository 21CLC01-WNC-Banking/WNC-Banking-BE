package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	AddNewAccount(ctx *gin.Context, customerId int64) error
	InternalTransfer(ctx *gin.Context, transferReq model.InternalTransferRequest) error
	GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.User, error)
}
