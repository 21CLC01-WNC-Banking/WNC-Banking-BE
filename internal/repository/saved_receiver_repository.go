package repository

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type SavedReceiverRepository interface {
	CreateCommand(ctx context.Context, savedReceiver *entity.SavedReceiver) error
}
