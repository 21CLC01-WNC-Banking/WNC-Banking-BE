package entity

import "time"

type SavedReceiver struct {
	ID                    int64      `db:"id" json:"id,omitempty"`
	CustomerId            int64      `db:"customer_id" json:"customerId,omitempty"`
	ReceiverAccountNumber string     `db:"receiver_account_number" json:"receiverAccountNumber,omitempty"`
	ReceiverNickname      string     `db:"receiver_nickname" json:"receiverNickname,omitempty"`
	BankId                *int64     `db:"bank_id" json:"bankId,omitempty"`
	CreatedAt             *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt             *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
