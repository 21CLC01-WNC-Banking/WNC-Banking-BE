package entity

import "time"

type User struct {
	Id          int64      `db:"id" json:"id,omitempty"`
	Email       string     `db:"email" json:"email,omitempty"`
	Name        string     `db:"name" json:"name,omitempty"`
	RoleId      int64      `db:"role_id" json:"role_id,omitempty"`
	PhoneNumber string     `db:"phone_number" json:"phoneNumber,omitempty"`
	Password    string     `db:"password" json:"password,omitempty"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
