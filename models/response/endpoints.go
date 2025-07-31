package response

import (
	"time"

	"github.com/nam2184/mymy/models/db"
)

type PostAuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"exp"`
	User         db.User   `json:"user"`
}

type PostSignUpResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type GetRefreshAuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"exp"`
	User         db.User   `json:"user"`
}

type GetChats struct {
	Chats []db.Chat `json:"chats"`
	Users []db.User `json:"users"`
}

type GetEncryptedMessages struct {
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

func (g GetEncryptedMessages) Id() interface{} {
	return g.ID
}

func (g GetEncryptedMessages) TableName() string {
	return "encrypted_messages"
}

type GetMessages struct {
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

type GetUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type GetCount struct {
	Count int64 `json:"count"`
}
