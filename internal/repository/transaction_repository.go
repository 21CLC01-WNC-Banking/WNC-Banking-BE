package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type TransactionRepository interface {
	CreateCommand(ctx context.Context, transaction *entity.Transaction) (string, error)
	GetTransactionBySourceNumberAndIdQuery(ctx context.Context, sourceNumber string, id string) (*entity.Transaction, error)
	UpdateStatusCommand(ctx context.Context, transaction *entity.Transaction) error
	UpdateBalancesCommand(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionByAccountNumber(ctx context.Context, accountNumber string) ([]entity.Transaction, error)
	GetTransactionByIdQuery(ctx context.Context, id string) (*entity.Transaction, error)
	GetReceivedDebtReminderByCustomerIdQuery(ctx context.Context, customerId int64) (*[]entity.Transaction, error)
	GetSentDebtReminderByCustomerIdQuery(ctx context.Context, customerId int64) (*[]entity.Transaction, error)
	GetTransactionByAccountNumberAndIdQuery(ctx context.Context, accountNumber string, id string) (*entity.Transaction, error)
	GetExternalTransactionsWithFilter(ctx context.Context, fromDate, toDate string, bankId int64) ([]entity.Transaction, error)
}
