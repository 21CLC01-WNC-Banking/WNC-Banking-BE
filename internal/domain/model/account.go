package model

import "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"

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
