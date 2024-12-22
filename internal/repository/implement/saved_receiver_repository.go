package repositoryimplement

import (
	"context"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type SavedReceiverRepository struct {
	db *sqlx.DB
}

func NewSavedReceiverRepository(db database.Db) repository.SavedReceiverRepository {
	return &SavedReceiverRepository{
		db: db,
	}
}

func (repo *SavedReceiverRepository) CreateCommand(ctx context.Context, savedReceiver *entity.SavedReceiver) error {
	insertQuery := `INSERT INTO saved_receivers(customer_id, receiver_account_number, receiver_nickname, bank_id) VALUES (:customer_id, :receiver_account_number, :receiver_nickname, :bank_id)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, savedReceiver)
	if err != nil {
		return err
	}
	return nil
}
