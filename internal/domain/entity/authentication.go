package entity

import "time"

type Authentication struct {
	ID        	int64      `db:"id" json:"id"`
	CustomerID  int64     `db:"customer_id" json:"customerId"`
	RefreshToken string 	`db:"refresh_token" json:"refreshToken"`
	CreatedAt 	*time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt 	*time.Time `db:"updated_at" json:"updatedAt"`
	DeleteddAt 	*time.Time `db:"deleted_at" json:"deleteddAt"`
}