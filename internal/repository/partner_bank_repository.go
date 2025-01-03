package repository

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type PartnerBankRepository interface {
	CreateCommand(ctx context.Context, partnerBank *entity.PartnerBank) error
	GetOneByBankCode(ctx context.Context, bankCode string) (*entity.PartnerBank, error)
	GetListBank(ctx context.Context) ([]entity.PartnerBank, error)
	GetBankById(ctx context.Context, bankId int64) (*entity.PartnerBank, error)
}
