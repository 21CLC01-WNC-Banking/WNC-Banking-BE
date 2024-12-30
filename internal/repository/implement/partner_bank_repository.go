package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type PartnerBankRepository struct {
	db *sqlx.DB
}

func NewPartnerBankRepository(db database.Db) repository.PartnerBankRepository {
	return &PartnerBankRepository{db: db}
}

func (repo *PartnerBankRepository) CreateCommand(ctx context.Context, partnerBank *entity.PartnerBank) error {
	insertQuery := `INSERT INTO partner_banks(bank_code,bank_name,short_name,logo_url,research_api,transfer_api)
				VALUES (:bank_code,:bank_name,:short_name,:logo_url,:research_api,:transfer_api)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, partnerBank)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PartnerBankRepository) GetOneByBankCode(ctx context.Context, bankCode string) (*entity.PartnerBank, error) {
	query := `SELECT id,bank_code,bank_name,short_name,logo_url,research_api,transfer_api
				FROM partner_banks WHERE bank_code=?`
	var partnerBank entity.PartnerBank
	err := repo.db.QueryRowxContext(ctx, query, bankCode).StructScan(&partnerBank)
	if err != nil {
		return nil, err
	}
	return &partnerBank, nil
}
