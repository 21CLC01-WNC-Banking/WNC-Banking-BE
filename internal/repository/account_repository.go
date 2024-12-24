package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type AccountRepository interface {
	CreateCommand(ctx context.Context, account *entity.Account) error
	UpdateBalanceCommand(ctx context.Context, number string, amount int64) (int64, error)
	GetOneByNumberQuery(ctx context.Context, number string) (*entity.Account, error)
	GetOneByCustomerIdQuery(ctx context.Context, customerId int64) (*entity.Account, error)
}
