package repositoryimplement

import (
	"context"
	"fmt"

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

func (repo *SavedReceiverRepository) ExistsByAccountNumberAndBankID(ctx context.Context, accountNumber string, bankID *int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM saved_receivers
			WHERE receiver_account_number = ? AND (bank_id = ? OR (? IS NULL AND bank_id IS NULL))
		)
	`
	err := repo.db.QueryRowContext(ctx, query, accountNumber, bankID, bankID).Scan(&exists)
	return exists, err
}

func (repo *SavedReceiverRepository) GetAllQuery(ctx context.Context, userId int64) (*[]entity.SavedReceiver, error) {
	var receivers []entity.SavedReceiver
	query := `SELECT * FROM saved_receivers WHERE customer_id = ?`
	err := repo.db.SelectContext(ctx, &receivers, query, userId)
	if err != nil {
		return nil, err
	}
	return &receivers, nil
}

func (repo *SavedReceiverRepository) UpdateNameByIdQuery(ctx context.Context, id int64, userId int64, newNickname string) error {
	query := `
		UPDATE saved_receivers
		SET receiver_nickname = ?
		WHERE id = ? AND customer_id = ?
	`
	result, err := repo.db.ExecContext(ctx, query, newNickname, id, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no receiver found")
	}

	return nil
}

func (repo *SavedReceiverRepository) DeleteReceiverByIdQuery(ctx context.Context, id int64, userId int64) error {
	query := `
		DELETE FROM saved_receivers
		WHERE id = ? AND customer_id = ?
	`
	result, err := repo.db.ExecContext(ctx, query, id, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no receiver found")
	}

	return nil
}
