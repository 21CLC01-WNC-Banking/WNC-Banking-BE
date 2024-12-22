package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type TransactionRepository interface {
	CreateCommand(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionBySourceNumberAndIdQuery(ctx context.Context, sourceNumber string, id string) (*entity.Transaction, error)
	UpdateStatusCommand(ctx context.Context, transaction *entity.Transaction) error
}
