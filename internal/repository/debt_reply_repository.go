package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type DebtReplyRepository interface {
	CreateCommand(ctx context.Context, reply *entity.DebtReply) error
	GetReplyByDebtIdQuery(ctx context.Context, debtId string) (*entity.DebtReply, error)
}
