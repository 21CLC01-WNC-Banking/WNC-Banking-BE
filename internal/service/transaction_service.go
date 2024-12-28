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
	AddDebtReminder(ctx *gin.Context, debtReminder model.DebtReminderRequest) error
	CancelDebtReminder(ctx *gin.Context, debtReminderId string, debtReply model.DebtReminderReplyRequest) error
	GetReceivedDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error)
	GetSentDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error)
	GetTransactionsByCustomerId(ctx *gin.Context, customerId int64) ([]model.GetTransactionsResponse, error)
	GetTransactionByIdAndCustomerId(ctx *gin.Context, customerId int64, id string) (*model.GetTransactionsResponse, error)
}
