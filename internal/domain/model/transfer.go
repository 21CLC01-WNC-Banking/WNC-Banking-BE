package model

type PreInternalTransferRequest struct {
	SourceAccountNumber string `json:"sourceAccountNumber" binding:"required"`
	TargetAccountNumber string `json:"targetAccountNumber" binding:"required"`
	Amount              int64  `json:"amount" binding:"required,min=0"`
	IsSourceFee         *bool  `json:"isSourceFee" binding:"required"`
	Description         string `json:"description" binding:"required"`
	Type                string `json:"type" binding:"required"`
}

type InternalTransferRequest struct {
	TransactionId string `json:"transactionId" binding:"required"`
	Otp           string `json:"otp" binding:"required"`
}

type UserInternalTransferResponse struct {
	TransactionId string `json:"transactionId" binding:"required"`
}
