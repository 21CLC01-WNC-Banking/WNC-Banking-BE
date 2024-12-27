package repositoryimplement

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db database.Db) repository.NotificationRepository {
	return &NotificationRepository{db: db}
}

func (n *NotificationRepository) GetAllByUserId(ctx context.Context, userId int64) ([]entity.Notification, error) {
	query := `
		SELECT 
			id, type, title, content, is_seen, user_id, created_at, updated_at, deleted_at
		FROM 
			notifications
		WHERE 
			user_id = ?
		ORDER BY 
			created_at DESC
	`

	var notifications []entity.Notification
	err := n.db.SelectContext(ctx, &notifications, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.Notification{}, nil
		}
		return nil, fmt.Errorf("failed to fetch notifications for user %d: %w", userId, err)
	}

	return notifications, nil
}

func (n *NotificationRepository) PatchSeenById(ctx context.Context, id int64) error {
	query := `UPDATE notifications SET is_seen = true WHERE id = ?`

	_, err := n.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func (n *NotificationRepository) GetCustomerIdById(ctx context.Context, id int64) int64 {
	query := `
		SELECT 
			user_id
		FROM 
			notifications
		WHERE 
			id = ?
	`

	var userId int64

	err := n.db.GetContext(ctx, &userId, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}

	return userId
}
