package model

import "time"

type AddAmountToAccountRequest struct {
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}

type GetTransactionsResponse struct {
	Id          string     `json:"id"`
	Amount      int64      `json:"amount"`
	CreatedAt   *time.Time `json:"createdAt"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Balance     int64      `json:"balance"`
}

type GetTransactionsByCustomerResponse struct {
	CustomerName string                    `json:"customerName"`
	Transactions []GetTransactionsResponse `json:"transactions"`
}
