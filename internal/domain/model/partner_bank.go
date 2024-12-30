package model

type PartnerBankRequest struct {
	BankCode    string `json:"bankCode"`
	BankName    string `json:"bankName"`
	ShortName   string `json:"shortName"`
	LogoUrl     string `json:"logoUrl"`
	ResearchApi string `json:"researchApi"`
	TransferApi string `json:"transferApi"`
}
