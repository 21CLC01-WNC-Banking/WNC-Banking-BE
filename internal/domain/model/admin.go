package model

import "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"

type CreateStaffRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required, min=8,max=255"`
}

type CreateStaffResponse struct {
	Id int64 `json:"id"`
}

type UpdateStaffRequest struct {
	Id          int64  `json:"id" binding:"required"` // required
	Name        string `json:"name"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type GetExternalTransactionRequest struct {
	BankId   int64  `json:"bankId"`
	FromDate string `json:"fromDate" binding:"required"` // yyyy-mm-dd
	ToDate   string `json:"toDate" binding:"required"`   // yyyy-mm-dd
}

type GetExternalTransactionResponse struct {
	entity.Transaction
	PartnerBankId        int64  `json:"partnerBankId" db:"bank_id"`
	PartnerBankShortName string `json:"partnerBankShortName" db:"bank_short_name"`
	PartnerBankName      string `json:"partnerBankName" db:"bank_name"`
	PartnerBankCode      string `json:"partnerBankCode" db:"bank_code"`
}
