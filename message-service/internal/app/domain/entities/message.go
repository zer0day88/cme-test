package entities

import "time"

type Message struct {
	ID          string    `json:"id"`
	Sender      string    `json:"sender"`
	Recipient   string    `json:"recipient"`
	Content     string    `json:"content"`
	MessageTime time.Time `json:"message_time"`
}
