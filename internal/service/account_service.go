package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	AddNewAccount(ctx *gin.Context, customerId int64) error
	GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.Customer, error)
	UpdateBalanceByAccountNumber(ctx *gin.Context, balance int64, number string) error
	GetAccountByCustomerId(ctx *gin.Context, customerId int64) (*entity.Account, error)
	GetAccountByNumber(ctx *gin.Context, number string) (*entity.Account, error)
}
