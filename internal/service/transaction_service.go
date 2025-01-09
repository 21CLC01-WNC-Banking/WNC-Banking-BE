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
	GetTransactionsByCustomerId(ctx *gin.Context, customerId int64) ([]model.GetTransactionsResponseSum, error)
	GetTransactionByIdAndCustomerId(ctx *gin.Context, customerId int64, id string) (*model.GetTransactionsResponse, error)
	PreDebtTransfer(ctx *gin.Context, transferReq model.PreDebtTransferRequest) error
	ReceiveExternalTransfer(ctx *gin.Context, transferReq model.ExternalPayload, partnerBankId int64) error
	PreExternalTransfer(ctx *gin.Context, transferReq model.PreExternalTransferRequest) (string, error)
	ExternalTransfer(ctx *gin.Context, transferReq model.TransferRequest) (*entity.Transaction, error)
}

func TransactionUtilsEntityToResponse(transaction entity.Transaction, sourceAccountNumber string) model.GetTransactionsResponse {
	var amount int64
	var balance int64

	if transaction.TargetAccountNumber == sourceAccountNumber {
		amount = transaction.Amount
		balance = transaction.TargetBalance
	} else {
		amount = transaction.Amount * -1
		balance = transaction.SourceBalance
	}

	return model.GetTransactionsResponse{
		Id:                  transaction.Id,
		Amount:              amount,
		CreatedAt:           transaction.CreatedAt,
		Description:         transaction.Description,
		Type:                transaction.Type,
		Balance:             balance,
		SourceAccountNumber: transaction.SourceAccountNumber,
		TargetAccountNumber: transaction.TargetAccountNumber,
	}
}
