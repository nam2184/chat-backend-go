package db

import (
	"time"
)

const MessageTableName = "messages"

// Message represents an individual message in a chat
type Message struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id"`                  // Primary key
	ChatID     uint      `gorm:"not null;index;column:chat_id" db:"chat_id" json:"chat_id"`             // Foreign key reference to the chat
	SenderID   uint      `gorm:"not null;index;column:sender_id" db:"sender_id" json:"sender_id"`       // ID of the sender
	SenderName string    `gorm:"not null;index;column:sender_name" db:"sender_name" json:"sender_name"` // ID of the sender
	ReceiverID uint      `gorm:"not null;index;column:receiver_id" db:"receiver_id" json:"receiver_id"` // ID of the second user
	Content    string    `gorm:"type:text;not null;column:content" db:"content" json:"content"`         // Message content
	Image      []byte    `gorm:"type:blob;column:image" json:"image"`                                   // Store as binary blob
	Type       string    `gorm:"type:text;not null;column:type" db:"type" json:"type"`
	IsTyping   bool      `gorm:"column:is_typing" db:"is_typing" json:"is_typing"`
	Timestamp  time.Time `gorm:"autoCreateTime;column:timestamp" db:"timestamp" json:"timestamp"` // Auto-managed timestamp for creation
}

func (m Message) TableName() string {
	return MessageTableName
}

func (m Message) Id() interface{} {
	return m.ID
}
