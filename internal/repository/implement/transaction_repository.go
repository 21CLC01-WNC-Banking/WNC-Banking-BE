package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db database.Db) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateCommand(ctx context.Context, transaction *entity.Transaction) error {
	//insert new transaction
	insertQuery := `INSERT INTO transactions(id, source_account_number, target_account_number,
											amount, bank_id, type, description, status, is_source_fee,
                         					source_balance, target_balance) VALUES
											(:id, :source_account_number, :target_account_number,
											 :amount, :bank_id, :type, :description, :status,
											 :is_source_fee, :source_balance, :target_balance)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, transaction)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepository) UpdateStatusCommand(ctx context.Context, transaction *entity.Transaction) error {
	query := `UPDATE transactions SET status = :status WHERE id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, transaction)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepository) GetTransactionBySourceNumberAndIdQuery(ctx context.Context, sourceNumber string, id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	query := "SELECT * FROM transactions WHERE source_account_number = ? AND id = ?"
	err := repo.db.QueryRowxContext(ctx, query, sourceNumber, id).StructScan(&transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
