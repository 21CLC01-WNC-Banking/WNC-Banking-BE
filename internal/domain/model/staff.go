package model

import "time"

type AddAmountToAccountRequest struct {
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Message       string `json:"message"`
}

type GetTransactionsResponse struct {
	Id                  string     `json:"id"`
	Amount              int64      `json:"amount"`
	CreatedAt           *time.Time `json:"createdAt"`
	Description         string     `json:"description"`
	Type                string     `json:"type"`
	Balance             int64      `json:"balance"`
	SourceAccountNumber string     `json:"sourceAccountNumber"`
	TargetAccountNumber string     `json:"targetAccountNumber"`
}

type GetTransactionsByCustomerResponse struct {
	CustomerName string                    `json:"customerName"`
	Transactions []GetTransactionsResponse `json:"transactions"`
}
