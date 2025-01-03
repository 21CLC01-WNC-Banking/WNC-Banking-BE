package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type AuthenticationRepository interface {
	CreateCommand(ctx context.Context, refreshToken entity.Authentication) error
	UpdateCommand(ctx context.Context, refreshToken entity.Authentication) error
	GetOneByCustomerIdQuery(ctx context.Context, customerId int64) (*entity.Authentication, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
}
