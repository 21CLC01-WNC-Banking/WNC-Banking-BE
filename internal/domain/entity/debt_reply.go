package entity

import "time"

type DebtReply struct {
	ID             int64      `db:"id" json:"id,omitempty"`
	DebtReminderId string     `db:"debt_reminder_id" json:"debtReminderId,omitempty"`
	UserReplyName  string     `db:"user_reply_name" json:"userReplyName,omitempty"`
	Content        string     `db:"content" json:"content,omitempty"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
}
