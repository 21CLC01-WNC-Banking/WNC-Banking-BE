package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type CustomerRepository interface {
	CreateCommand(ctx context.Context, customer *entity.User) error
	GetOneByEmailQuery(ctx context.Context, email string) (*entity.User, error)
	GetIdByEmailQuery(ctx context.Context, email string) (int64, error)
	UpdatePasswordByIdQuery(ctx context.Context, id int64, password string) error
	GetOneByIdQuery(ctx context.Context, id int64) (*entity.User, error)
	GetCustomerByAccountNumberQuery(ctx context.Context, number string) (*entity.User, error)
	DeleteById(ctx context.Context, userId int64) error
}
