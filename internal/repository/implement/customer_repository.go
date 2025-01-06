package repositoryimplement

import (
	"context"
	"fmt"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db database.Db) repository.CustomerRepository {
	return &CustomerRepository{db: db}
}

func (repo *CustomerRepository) CreateCommand(ctx context.Context, customer *entity.User) error {
	// Insert the new customer
	insertQuery := `INSERT INTO users(email, name, role_id, phone_number, password) VALUES (:email, :name, :role_id, :phone_number, :password)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, customer)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CustomerRepository) GetOneByEmailQuery(ctx context.Context, email string) (*entity.User, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repo *CustomerRepository) GetOneByIdQuery(ctx context.Context, id int64) (*entity.User, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE id = ? AND users.deleted_at IS NULL"
	err := repo.db.QueryRowxContext(ctx, query, id).StructScan(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repo *CustomerRepository) GetIdByEmailQuery(ctx context.Context, email string) (int64, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&customer)
	if err != nil {
		return 0, err
	}
	return customer.Id, nil
}

func (repo *CustomerRepository) UpdatePasswordByIdQuery(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password = ? WHERE id = ? AND users.deleted_at IS NULL"
	_, err := repo.db.ExecContext(ctx, query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CustomerRepository) GetCustomerByAccountNumberQuery(ctx context.Context, number string) (*entity.User, error) {
	var customer entity.User
	query := `
				SELECT users.* FROM users 
				JOIN accounts ON users.id = accounts.customer_id AND accounts.number = ?
				WHERE users.deleted_at IS NULL
			 `
	err := repo.db.QueryRowxContext(ctx, query, number).StructScan(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// also soft delete user from saved_receivers and hard delete from authentications
func (repo *CustomerRepository) DeleteById(ctx context.Context, userId int64) error {
	// Start a transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is properly rolled back on error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Soft delete from users
	query := `
		UPDATE users 
		SET deleted_at = CURRENT_TIMESTAMP() 
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to soft delete user with ID %d: %w", userId, err)
	}

	// Soft delete from accounts
	query = `
		UPDATE accounts 
		SET deleted_at = CURRENT_TIMESTAMP() 
		WHERE customer_id = ? AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to soft delete user with ID %d: %w", userId, err)
	}

	// Soft delete from saved_receivers
	query = `
		UPDATE saved_receivers 
		SET deleted_at = CURRENT_TIMESTAMP() 
		WHERE customer_id = ? AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to soft delete saved_receivers with userID %d: %w", userId, err)
	}

	// Hard delete from authentications
	query = `
		DELETE FROM authentications 
		WHERE user_id = ?
	`
	_, err = tx.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to hard delete authentications with userID %d: %w", userId, err)
	}

	return nil
}
