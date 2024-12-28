package model

import "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"

type DebtReminderRequest struct {
	SourceAccountNumber string `json:"sourceAccountNumber" binding:"required"`
	TargetAccountNumber string `json:"targetAccountNumber" binding:"required"`
	Amount              int64  `json:"amount" binding:"required,min=0"`
	Description         string `json:"description" binding:"required"`
	Type                string `json:"type" binding:"required"`
}

type DebtReminderReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

type DebtReminderResponse struct {
	Sender       string              `json:"sender" binding:"required"`
	Receiver     string              `json:"receiver" binding:"required"`
	DebtReminder *entity.Transaction `json:"debtReminder" binding:"required"`
	Reply        *entity.DebtReply   `json:"reply"`
}
