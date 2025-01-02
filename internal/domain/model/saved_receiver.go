package model

type Receiver struct {
	BankId                int64  `json:"bankId"`
	ReceiverAccountNumber string `json:"receiverAccountNumber" binding:"required"`
	ReceiverNickname      string `json:"receiverNickname" binding:"required"`
}

type SavedReceiverResponse struct {
	ID                    int64  `json:"id" binding:"required"`
	ReceiverAccountNumber string `json:"receiverAccountNumber" binding:"required"`
	ReceiverNickname      string `json:"receiverNickname" binding:"required"`
	BankId                *int64 `json:"bankId"`
}

type UpdateReceiverRequest struct {
	NewNickname string `json:"newNickname" binding:"required"`
}
