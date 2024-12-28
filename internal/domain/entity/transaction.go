package entity

import "time"

type Transaction struct {
	Id                  string     `db:"id" json:"id"`
	SourceAccountNumber string     `db:"source_account_number" json:"sourceAccountNumber"`
	TargetAccountNumber string     `db:"target_account_number" json:"targetAccountNumber"`
	Amount              int64      `db:"amount" json:"amount"`
	BankId              *int64     `db:"bank_id" json:"bankId"`
	Type                string     `db:"type" json:"type" enum:"internal,external,payment"`
	Description         string     `db:"description" json:"description"`
	Status              string     `db:"status" json:"status" enum:"pending,failed,success"`
	IsSourceFee         *bool      `db:"is_source_fee" json:"isSourceFee"`
	SourceBalance       int64      `db:"source_balance" json:"sourceBalance"`
	TargetBalance       int64      `db:"target_balance" json:"targetBalance"`
	CreatedAt           *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deletedAt"`
}
