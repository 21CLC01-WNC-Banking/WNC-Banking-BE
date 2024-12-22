package entity

import "time"

type Authentication struct {
	ID           int64      `db:"id" json:"id"`
	UserId       int64      `db:"user_id" json:"userId"`
	RefreshToken string     `db:"refresh_token" json:"refreshToken"`
	CreatedAt    *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deletedAt"`
}
