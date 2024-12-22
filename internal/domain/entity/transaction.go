package entity

import "time"

type Transaction struct {
	Id                  string     `db:"id" json:"id"`
	SourceAccountNumber string     `db:"source_account_number" json:"source_account_number"`
	TargetAccountNumber string     `db:"target_account_number" json:"target_account_number"`
	Amount              int64      `db:"amount" json:"amount"`
	BankId              *int64     `db:"bank_id" json:"bank_id"`
	Type                string     `db:"type" json:"type"`
	Description         string     `db:"description" json:"description"`
	Status              string     `db:"status" json:"status"`
	IsSourceFee         *bool      `db:"is_source_fee" json:"is_source_fee"`
	SourceBalance       int64      `db:"source_balance" json:"source_balance"`
	TargetBalance       int64      `db:"target_balance" json:"target_balance"`
	CreatedAt           *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deletedAt"`
}
