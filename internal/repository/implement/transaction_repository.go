package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/generate_number_code"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db database.Db) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateCommand(ctx context.Context, transaction *entity.Transaction) (string, error) {
	transactionId := generate_number_code.GenerateRandomNumber(10)

	transaction.Id = transactionId
	//insert new transaction
	insertQuery := `INSERT INTO transactions(id, source_account_number, target_account_number,
											amount, bank_id, type, description, status, is_source_fee,
                         					source_balance, target_balance) VALUES
											(:id, :source_account_number, :target_account_number,
											 :amount, :bank_id, :type, :description, :status,
											 :is_source_fee, :source_balance, :target_balance)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, transaction)
	if err != nil {
		return "", err
	}
	return transactionId, err
}

func (repo *TransactionRepository) UpdateBalancesCommand(ctx context.Context, transaction *entity.Transaction) error {
	query := `UPDATE transactions 
			  SET source_balance = :source_balance,
			  target_balance = :target_balance
			  WHERE id = :id`

	_, err := repo.db.NamedExecContext(ctx, query, transaction)
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

func (repo *TransactionRepository) GetTransactionByAccountNumber(ctx context.Context, accountNumber string) ([]entity.Transaction, error) {
	query := `
		SELECT 
			id, source_account_number, target_account_number, amount, bank_id, 
			type, description, status, is_source_fee, source_balance, 
			target_balance, created_at, updated_at, deleted_at
		FROM transactions
		WHERE source_account_number = ? OR target_account_number = ? AND status = "success"
		ORDER BY created_at DESC
	`

	var transactions []entity.Transaction
	err := repo.db.SelectContext(ctx, &transactions, query, accountNumber, accountNumber)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *TransactionRepository) GetTransactionByIdQuery(ctx context.Context, id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	query := "SELECT * FROM transactions WHERE id = ?"
	err := repo.db.QueryRowxContext(ctx, query, id).StructScan(&transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
