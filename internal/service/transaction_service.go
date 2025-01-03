package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type TransactionService interface {
	PreInternalTransfer(ctx *gin.Context, transferReq model.PreInternalTransferRequest) (string, error)
	SendOTPToEmail(ctx *gin.Context, email string, transactionId string) error
	InternalTransfer(ctx *gin.Context, transferReq model.TransferRequest) (*entity.Transaction, error)
	AddDebtReminder(ctx *gin.Context, debtReminder model.DebtReminderRequest) error
	CancelDebtReminder(ctx *gin.Context, debtReminderId string, debtReply model.DebtReminderReplyRequest) error
	GetReceivedDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error)
	GetSentDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error)
	GetTransactionsByCustomerId(ctx *gin.Context, customerId int64) ([]model.GetTransactionsResponse, error)
	GetTransactionByIdAndCustomerId(ctx *gin.Context, customerId int64, id string) (*model.GetTransactionsResponse, error)
	PreDebtTransfer(ctx *gin.Context, transferReq model.PreDebtTransferRequest) error
	ReceiveExternalTransfer(ctx *gin.Context, transferReq model.ExternalTransactionRequest, partnerBankId int64) error
	PreExternalTransfer(ctx *gin.Context, transferReq model.PreExternalTransferRequest) (string, error)
	ExternalTransfer(ctx *gin.Context, transferReq model.TransferRequest) (*entity.Transaction, error)
}
