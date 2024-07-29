package model

import "time"

type Message struct {
	ID          string    `json:"id"`
	Sender      string    `json:"sender"`
	Recipient   string    `json:"recipient"`
	Content     string    `json:"content"`
	MessageTime time.Time `json:"message_time"`
}

type SendRequest struct {
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
}
