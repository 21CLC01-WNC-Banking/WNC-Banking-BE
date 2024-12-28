package entity

import "time"

type DebtReply struct {
	ID             int64      `db:"id" json:"id,omitempty"`
	DebtReminderId string     `db:"debt_reminder_id" json:"debtReminderId,omitempty"`
	UserReplyId    int64      `db:"user_reply_id" json:"userReplyId,omitempty"`
	Content        string     `db:"content" json:"content,omitempty"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
}
