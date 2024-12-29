package model

type NotificationRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	DeviceId int    `json:"deviceId"`
}
