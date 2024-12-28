package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type StaffRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetOneById(ctx context.Context, id int64) (*entity.User, error)
}
