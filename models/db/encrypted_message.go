package db

import (
	"time"
)

const EncryptedMessageTableName = "encrypted_messages"

// Message represents an individual message in a chat
type EncryptedMessage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id"`
	ChatID     uint      `gorm:"not null;index;column:chat_id" db:"chat_id" json:"chat_id"`
	SenderID   uint      `gorm:"not null;index;column:sender_id" db:"sender_id" json:"sender_id"`
	SenderName string    `gorm:"not null;index;column:sender_name" db:"sender_name" json:"sender_name"`
	ReceiverID uint      `gorm:"not null;index;column:receiver_id" db:"receiver_id" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null;column:content" db:"content" json:"content"` // Message content
	Image      []byte    `gorm:"type:bytea;column:image" db:"image" json:"image"`
	Type       string    `gorm:"type:text;not null;column:type" db:"type" json:"type"`
	IsTyping   bool      `gorm:"column:is_typing" db:"is_typing" json:"is_typing"`
	Timestamp  time.Time `gorm:"autoCreateTime;column:timestamp" db:"timestamp" json:"timestamp"`
}

func (m EncryptedMessage) TableName() string {
	return EncryptedMessageTableName
}

func (m EncryptedMessage) Id() interface{} {
	return m.ID
}
