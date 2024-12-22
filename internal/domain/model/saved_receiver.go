package model

type InternalReceiver struct {
	ReceiverAccountNumber string `json:"receiverAccountNumber" binding:"required"`
	ReceiverNickname      string `json:"receiverNickname" binding:"required"`
}

type ExternalReceiver struct {
	BankCode              int64  `json:"bankCode" binding:"required"`
	ReceiverAccountNumber string `json:"receiverAccountNumber" binding:"required"`
	ReceiverNickname      string `json:"receiverNickname" binding:"required"`
}
