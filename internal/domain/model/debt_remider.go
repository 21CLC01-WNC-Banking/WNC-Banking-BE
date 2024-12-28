package model

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
