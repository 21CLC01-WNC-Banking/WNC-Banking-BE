package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type SavedReceiverRepository interface {
	CreateCommand(ctx context.Context, savedReceiver *entity.SavedReceiver) error
	ExistsByAccountNumberAndBankID(ctx context.Context, accountNumber string, bankID *int64) (bool, error)
	GetAllQuery(ctx context.Context, userId int64) (*[]entity.SavedReceiver, error)
	UpdateNameByIdQuery(ctx context.Context, id int64, userId int64, newNickname string) error
	DeleteReceiverByIdQuery(ctx context.Context, id int64, userId int64) error
}
