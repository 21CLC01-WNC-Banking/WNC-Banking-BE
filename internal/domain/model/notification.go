package model

type NotificationRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}
