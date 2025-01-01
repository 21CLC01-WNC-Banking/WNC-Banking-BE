package model

import "time"

type NotificationRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	DeviceId int    `json:"deviceId"`
}

type TransactionNotificationContent struct {
	DeviceId      int        `json:"deviceId"`
	Name          string     `json:"name"`
	Amount        int        `json:"amount"`
	TransactionId string     `json:"transactionId"`
	Type          string     `json:"type"`
	CreatedAt     *time.Time `json:"createdAt"`
}
