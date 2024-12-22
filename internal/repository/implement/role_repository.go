package repositoryimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db database.Db) repository.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByUserId(c *gin.Context, userId int64) (*entity.Role, error) {
	// Query to get the role based on the userId
	query := `
		SELECT r.id, r.name
		FROM roles r
		JOIN users u ON u.role_id = r.id
		WHERE u.id = ?;
	`

	// Create a variable to hold the role
	var role entity.Role

	// Execute the query
	err := r.db.GetContext(c, &role, query, userId)
	if err != nil {
		return nil, err // Other errors
	}

	return &role, nil
}
