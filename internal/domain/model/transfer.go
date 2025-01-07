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

type ExternalPayload struct {
	SrcAccountNumber string `json:"srcAccountNumber" binding:"required"`
	SrcBankCode      string `json:"srcBankCode" binding:"required"`
	DesAccountNumber string `json:"desAccountNumber" binding:"required"`
	Amount           int64  `json:"amount" binding:"required,min=0"`
	Description      string `json:"description" binding:"required"`
	IsSourceFee      *bool  `json:"isSourceFee" binding:"required"`
	Exp              string `json:"exp" binding:"required"`
}

type ExternalTransferRequest struct {
	SrcAccountNumber string `json:"srcAccountNumber" binding:"required"`
	SrcBankCode      string `json:"srcBankCode" binding:"required"`
	DesAccountNumber string `json:"desAccountNumber" binding:"required"`
	Amount           int64  `json:"amount" binding:"required,min=0"`
	Description      string `json:"description" binding:"required"`
	IsSourceFee      *bool  `json:"isSourceFee" binding:"required"`
	Exp              string `json:"exp" binding:"required"`
	SignedData       string `json:"signedData" binding:"required"`
}

type ExternalTransactionResponse struct {
	Data       string `json:"data"`
	SignedData string `json:"signedData"`
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

type ExternalTransferToRSATeam struct {
	BankId               int64  `json:"bankId" binding:"required"`
	AccountNumber        string `json:"accountNumber" binding:"required"`
	ForeignAccountNumber string `json:"foreignAccountNumber" binding:"required"`
	Amount               int64  `json:"amount" binding:"required,min=0"`
	Description          string `json:"description" binding:"required"`
	Timestamp            int64  `json:"timestamp" binding:"required"`
}
