package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type NotificationRepository interface {
	GetAllByUserId(ctx context.Context, userId int64) ([]entity.Notification, error)
	PatchSeenById(ctx context.Context, id int64) error
	GetCustomerIdById(ctx context.Context, id int64) int64
}
