package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type StaffRepository struct {
	db *sqlx.DB
}

func NewStaffRepository(db database.Db) repository.StaffRepository {
	return &StaffRepository{db: db}
}

func (s *StaffRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	// SQL query to fetch users with the role name "staff"
	query := `
		SELECT u.id, u.email, u.name, u.role_id, u.phone_number, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE r.name = 'staff'
	`

	var users []entity.User
	err := s.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *StaffRepository) GetOneById(ctx context.Context, id int64) (*entity.User, error) {
	// SQL query to fetch a user with the role name "staff" by their ID
	query := `
		SELECT u.id, u.email, u.name, u.role_id, u.phone_number, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE r.name = 'staff' AND u.id = ?
	`

	var user entity.User
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *StaffRepository) CreateOne(ctx context.Context, staff *entity.User) (int64, error) {
	// SQL query to insert a new user with the role of "staff"
	query := `
		INSERT INTO users (email, name, role_id, phone_number, password)
		VALUES (:email, :name, (SELECT id FROM roles WHERE name = 'staff'), :phone_number, :password)
	`

	// Named query execution
	result, err := s.db.NamedExecContext(ctx, query, staff)
	if err != nil {
		return 0, err
	}

	// Fetch the last inserted ID
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}
