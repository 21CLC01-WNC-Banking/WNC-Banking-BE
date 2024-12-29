package model

type NotificationRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  int    `json:"userId"`
}
