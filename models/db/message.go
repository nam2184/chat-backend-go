package db

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/nam2184/mymy/models/body"
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

func ConvertToMessageDB(temp body.TempMessage) (Message, error) {
	var message Message
	if temp.Image != " " {
		base64Data := temp.Image

		if strings.HasPrefix(base64Data, "data:") {
			commaIndex := strings.Index(base64Data, ",")
			if commaIndex != -1 {
				base64Data = base64Data[commaIndex+1:] // Strip prefix
			}
		}
		image, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return Message{}, fmt.Errorf("base64 decode failed: %w", err)
		}
		if err != nil {
			return Message{}, err
		}
		message.Image = image
	}
	message.ID = temp.ID
	message.ChatID = temp.ChatID
	message.SenderID = temp.SenderID
	message.SenderName = temp.SenderName
	message.ReceiverID = temp.ReceiverID
	message.Content = temp.Content
	message.Type = temp.Type
	message.IsTyping = temp.IsTyping
	message.Timestamp = temp.Timestamp

	return message, nil
}
