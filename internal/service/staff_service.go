package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type StaffService interface {
	RegisterCustomer(ctx *gin.Context, customerRequest model.RegisterRequest) error
	AddAmountToAccount(ctx *gin.Context, request *model.AddAmountToAccountRequest) error
	GetTransactionsByAccountNumber(ctx *gin.Context, accountNumber string) (*model.GetTransactionsByCustomerResponse, error)
}
