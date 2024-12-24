package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	AddNewAccount(ctx *gin.Context, customerId int64) error
	GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.User, error)
	UpdateBalanceByAccountNumber(ctx *gin.Context, amount int64, number string) (int64, error)
	GetAccountByCustomerId(ctx *gin.Context, customerId int64) (*entity.Account, error)
	GetAccountByNumber(ctx *gin.Context, number string) (*entity.Account, error)
}
