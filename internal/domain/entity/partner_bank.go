package entity

import "time"

type PartnerBank struct {
	ID          int64      `db:"id" json:"id,omitempty"`
	BankCode    string     `db:"bank_code" json:"bankCode,omitempty"`
	BankName    string     `db:"bank_name" json:"bankName,omitempty"`
	ShortName   string     `db:"short_name" json:"shortName,omitempty"`
	LogoUrl     string     `db:"logo_url" json:"logoUrl,omitempty"`
	ResearchApi string     `db:"research_api" json:"researchApi,omitempty"`
	TransferApi string     `db:"transfer_api," json:"transferApi,omitempty"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}
