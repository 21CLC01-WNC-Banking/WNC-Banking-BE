package model

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
)

type SearchInternalAccountRequest struct {
	AccountNumber string `json:"accountNumber" binding:"required"`
}

type InternalAccountResponse struct {
	CustomerName  string `json:"customerName" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}

type AccountResponse struct {
	Account *entity.Account `json:"account" binding:"required"`
	Name    string          `json:"name" binding:"required"`
}

type AccountNumberInfoRequest struct {
	SrcBankCode      string `json:"srcBankCode"`
	DesAccountNumber string `json:"desAccountNumber"`
	Exp              string `json:"exp"`
}

type AccountNumberInfoResponse struct {
	DesAccountNumber string `json:"desAccountNumber"`
	DesAccountName   string `json:"desAccountName"`
}

type SearchExternalAccountRequest struct {
	BankId        int64  `json:"bankId" binding:"required"`
	TimeStamp     string `json:"timestamp" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}
