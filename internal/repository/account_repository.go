package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type AccountRepository interface {
	CreateCommand(ctx context.Context, account *entity.Account) error
	UpdateCommand(ctx context.Context, account entity.Account) error
	GetOneByNumberQuery(ctx context.Context, number string) (*entity.Account, error)
	GetOneByCustomerIdQuery(ctx context.Context, customerId int64) (*entity.Account, error)
}
