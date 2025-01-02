package model

type PartnerBankRequest struct {
	BankCode    string `json:"bankCode"`
	BankName    string `json:"bankName"`
	ShortName   string `json:"shortName"`
	LogoUrl     string `json:"logoUrl"`
	ResearchApi string `json:"researchApi"`
	TransferApi string `json:"transferApi"`
	PublicKey   string `json:"publicKey"`
}

type GetExternalAccountNameRequest struct {
	BankId        int64  `json:"bankId"`
	AccountNumber string `json:"accountNumber"`
}
