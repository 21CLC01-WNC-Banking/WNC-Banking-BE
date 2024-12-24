package model

type AddAmountToAccountRequest struct {
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}
