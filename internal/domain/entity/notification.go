package entity

import "time"

type Notification struct {
	Id        int        `db:"id" json:"id"`
	Type      string     `db:"type" json:"type"` // incoming_transfer, outgoing_transfer, debt_reminder
	Title     string     `db:"title" json:"title"`
	Content   string     `db:"content" json:"content"`
	IsSeen    bool       `db:"is_seen" json:"isSeen"`
	UserID    int        `db:"user_id" json:"userId"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}
