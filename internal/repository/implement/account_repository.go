package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db database.Db) repository.AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (repo *AccountRepository) CreateCommand(ctx context.Context, account *entity.Account) error {
	insertQuery := `INSERT INTO accounts(customer_id, number, balance) VALUES (:customer_id, :number, :balance)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, account)
	if err != nil {
		return err
	}
	return nil
}

func (repo *AccountRepository) UpdateBalanceCommand(ctx context.Context, number string, amount int64) (int64, error) {
	query := `
	UPDATE accounts
	SET balance = balance + ?
	WHERE number = ?
	`

	_, err := repo.db.ExecContext(ctx, query, amount, number)
	if err != nil {
		return 0, err
	}
	var newBalance int64
	selectQuery := `
	SELECT balance
	FROM accounts
	WHERE number = ?
	`
	err = repo.db.GetContext(ctx, &newBalance, selectQuery, number)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (repo *AccountRepository) GetOneByNumberQuery(ctx context.Context, number string) (*entity.Account, error) {
	var account entity.Account
	query := "SELECT * FROM accounts WHERE number = ?"
	err := repo.db.QueryRowxContext(ctx, query, number).StructScan(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *AccountRepository) GetOneByCustomerIdQuery(ctx context.Context, customerId int64) (*entity.Account, error) {
	var account entity.Account
	query := "SELECT * FROM accounts WHERE customer_id = ?"
	err := repo.db.QueryRowxContext(ctx, query, customerId).StructScan(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
