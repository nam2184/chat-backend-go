package request

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/nam2184/mymy/models/db"
)

type WSEncryptedMessage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id"`
	ChatID     uint      `gorm:"not null;index;column:chat_id" db:"chat_id" json:"chat_id"`
	SenderID   uint      `gorm:"not null;index;column:sender_id" db:"sender_id" json:"sender_id"`
	SenderName string    `gorm:"not null;index;column:sender_name" db:"sender_name" json:"sender_name"`
	ReceiverID uint      `gorm:"not null;index;column:receiver_id" db:"receiver_id" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null;column:content" db:"content" json:"content"` // Message content
	Image      string    `gorm:"type:bytea;column:image" db:"image" json:"image"`
	Type       string    `gorm:"type:text;not null;column:type" db:"type" json:"type"`
	IsTyping   bool      `gorm:"column:is_typing" db:"is_typing" json:"is_typing"`
	Timestamp  time.Time `gorm:"autoCreateTime;column:timestamp" db:"timestamp" json:"timestamp"`
}

func (temp WSEncryptedMessage) ConvertToEncryptedMessageDB() (db.EncryptedMessage, error) {
	var message db.EncryptedMessage
	if temp.Image != " " {
		base64Image := temp.Image
		if strings.HasPrefix(base64Image, "data:") {
			commaIndex := strings.Index(base64Image, ",")
			if commaIndex != -1 {
				base64Image = base64Image[commaIndex+1:] // Strip prefix
			}
		}
		image, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			return db.EncryptedMessage{}, fmt.Errorf("base64 decode failed: %w", err)
		}
		message.Image = image
	}
	message.ID = temp.ID
	message.ChatID = temp.ChatID
	message.SenderID = temp.SenderID
	message.SenderName = temp.SenderName
	message.ReceiverID = temp.ReceiverID
	message.Type = temp.Type
	message.IsTyping = temp.IsTyping
	message.Timestamp = temp.Timestamp
	message.Content = temp.Content
	return message, nil
}

type WSMessage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id"`
	ChatID     uint      `gorm:"not null;index;column:chat_id" db:"chat_id" json:"chat_id"`
	SenderID   uint      `gorm:"not null;index;column:sender_id" db:"sender_id" json:"sender_id"`
	SenderName string    `gorm:"not null;index;column:sender_name" db:"sender_name" json:"sender_name"`
	ReceiverID uint      `gorm:"not null;index;column:receiver_id" db:"receiver_id" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null;column:content" db:"content" json:"content"` // Message content
	Image      string    `gorm:"type:bytea;column:image" db:"image" json:"image"`
	Type       string    `gorm:"type:text;not null;column:type" db:"type" json:"type"`
	IsTyping   bool      `gorm:"column:is_typing" db:"is_typing" json:"is_typing"`
	Timestamp  time.Time `gorm:"autoCreateTime;column:timestamp" db:"timestamp" json:"timestamp"`
}

func (temp WSMessage) ConvertToMessageDB() (db.Message, error) {
	var message db.Message
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
			return db.Message{}, fmt.Errorf("base64 decode failed: %w", err)
		}
		if err != nil {
			return db.Message{}, err
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

type PostAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostAuthSignUp struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}

func (s PostAuthSignUp) ToAuth() db.Auth {
	return db.Auth{
		Username: s.Username,
		Password: s.Password,
	}
}

func (s PostAuthSignUp) ToUser() db.User {
	return db.User{
		Username:  s.Username,
		FirstName: s.FirstName,
		Surname:   s.Surname,
		Email:     s.Email,
		CreatedAt: time.Now(),
	}
}
