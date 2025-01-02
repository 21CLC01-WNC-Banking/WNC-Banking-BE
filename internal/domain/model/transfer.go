package model

type PreInternalTransferRequest struct {
	SourceAccountNumber string `json:"sourceAccountNumber" binding:"required"`
	TargetAccountNumber string `json:"targetAccountNumber" binding:"required"`
	Amount              int64  `json:"amount" binding:"required,min=0"`
	IsSourceFee         *bool  `json:"isSourceFee" binding:"required"`
	Description         string `json:"description" binding:"required"`
	Type                string `json:"type" binding:"required"`
}

type TransferRequest struct {
	TransactionId string `json:"transactionId" binding:"required"`
	Otp           string `json:"otp" binding:"required"`
}

type PreDebtTransferRequest struct {
	TransactionId string `json:"transactionId" binding:"required"`
}

type ExternalTransactionRequest struct {
	SrcAccountNumber string `json:"srcAccountNumber" binding:"required"`
	SrcBankCode      string `json:"srcBankCode" binding:"required"`
	DesAccountNumber string `json:"desAccountNumber" binding:"required"`
	Amount           int64  `json:"amount" binding:"required,min=0"`
	Description      string `json:"description" binding:"required"`
	IsSourceFee      *bool  `json:"isSourceFee" binding:"required"`
}

type ExternalTransactionData struct {
	SrcAccountNumber string `json:"srcAccountNumber" binding:"required"`
	SrcBankCode      string `json:"srcBankCode" binding:"required"`
	DesAccountNumber string `json:"desAccountNumber" binding:"required"`
	Amount           int64  `json:"amount" binding:"required,min=0"`
	Description      string `json:"description" binding:"required"`
	IsSourceFee      *bool  `json:"isSourceFee" binding:"required"`
	SignedData       string `json:"signedData" binding:"required"`
}

type ExternalTransactionResponse struct {
	Reply string `json:"reply"`
}

type PreExternalTransferRequest struct {
	SourceAccountNumber string `json:"sourceAccountNumber" binding:"required"`
	TargetAccountNumber string `json:"targetAccountNumber" binding:"required"`
	PartnerBankId       int64  `json:"partnerBankId" binding:"required"`
	Amount              int64  `json:"amount" binding:"required,min=0"`
	IsSourceFee         *bool  `json:"isSourceFee" binding:"required"`
	Description         string `json:"description" binding:"required"`
	Type                string `json:"type" binding:"required"`
}
