package body

import "time"

type TempMessage struct {
	ID         uint      `json:"id"`
	ChatID     uint      `json:"chat_id"`
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	Image      string    `json:"image"`
	Type       string    `json:"type"`
	IsTyping   bool      `json:"is_typing"`
	Timestamp  time.Time `json:"timestamp"`
}
