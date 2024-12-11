package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type CustomerRepository interface {
	CreateCommand(ctx context.Context, customer *entity.Customer) error
	GetOneByEmailQuery(ctx context.Context, email string) (*entity.Customer, error)
	GetMailByIdQuery(ctx context.Context, id int64) (string, error)
	UpdatePasswordByIdQuery(ctx context.Context, id int64, password string) error
}
