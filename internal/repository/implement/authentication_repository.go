package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AuthenticationRepository struct {
	db *sqlx.DB
}

func NewAuthenticationRepository(db database.Db) repository.AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func (repo *AuthenticationRepository) CreateCommand(ctx context.Context, authentication entity.Authentication) error {
	query := `
		INSERT INTO authentications (customer_id, refresh_token)
		VALUES (:customer_id, :refresh_token)
	`
	_, err := repo.db.NamedExecContext(ctx, query, authentication)
	return err
}

func (repo *AuthenticationRepository) UpdateCommand(ctx context.Context, authentication entity.Authentication) error {
	query := `
		UPDATE authentications
		SET refresh_token = :refresh_token
		WHERE customer_id = :customer_id
	`
	_, err := repo.db.NamedExecContext(ctx, query, authentication)
	return err
}

// find refresh token
func (repo *AuthenticationRepository) GetOneByCustomerIdQuery(ctx context.Context, customerId int64) (*entity.Authentication, error) {
	var authentication entity.Authentication
	query := `
		SELECT customer_id, refresh_token, created_at
		FROM authentications
		WHERE customer_id = ?
	`
	err := repo.db.GetContext(ctx, &authentication, query, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}
	return &authentication, nil
}
