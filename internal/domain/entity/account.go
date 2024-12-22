package entity

import "time"

type Account struct {
	ID         int64      `db:"id" json:"id,omitempty"`
	CustomerID int64      `db:"customer_id" json:"customerId,omitempty"`
	Number     string     `db:"number" json:"number,omitempty"`
	Balance    int64      `db:"balance" json:"balance,omitempty"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
