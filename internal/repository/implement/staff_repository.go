package repositoryimplement

import (
	"context"
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
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
		WHERE r.name = 'staff' AND u.deleted_at IS NULL
	`

	var users []entity.User
	err := s.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *StaffRepository) GetOneById(ctx context.Context, id int64) (*entity.User, error) {
	// SQL query to fetch a user with the role name "staff" by their Id
	query := `
		SELECT u.id, u.email, u.name, u.role_id, u.phone_number, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE r.name = 'staff' AND u.id = ? AND u.deleted_at IS NULL
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

	// Fetch the last inserted Id
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *StaffRepository) DeleteOne(ctx context.Context, id int64) error {
	query := `
		UPDATE users u
		SET deleted_at = CURRENT_TIMESTAMP()
		WHERE u.id = :id AND u.deleted_at IS NULL AND u.role_id = (SELECT id FROM roles WHERE name = 'staff')
	`

	result, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(httpcommon.ErrorMessage.SqlxNoRow)
	}

	return nil
}

func (s *StaffRepository) UpdateOneStaff(ctx context.Context, staff *entity.User) error {
	query := `UPDATE users SET `
	args := make(map[string]interface{})

	if staff.Name != "" {
		query += "name = :name, "
		args["name"] = staff.Name
	}
	if staff.Email != "" {
		query += "email = :email, "
		args["email"] = staff.Email
	}
	if staff.PhoneNumber != "" {
		query += "phone_number = :phone_number, "
		args["phone_number"] = staff.PhoneNumber
	}
	if staff.Password != "" {
		query += "password = :password, "
		args["password"] = staff.Password
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]

	query += " WHERE id = :id AND deleted_at IS NULL AND role_id = (SELECT id FROM roles WHERE name = 'staff')"
	args["id"] = staff.Id

	result, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(httpcommon.ErrorMessage.SqlxNoRow)
	}

	return nil
}
