package entity

import "time"

type DebtReminder struct {
	ID                  int64      `db:"id" json:"id,omitempty"`
	SourceAccountNumber string     `db:"source_account_number" json:"source_account_number"`
	TargetAccountNumber string     `db:"target_account_number" json:"target_account_number"`
	Amount              int64      `db:"amount" json:"amount"`
	Content             string     `db:"content" json:"content"`
	Status              string     `db:"status" json:"status" enum:"sent,paid,cancelled"`
	CreatedAt           *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deletedAt"`
}
