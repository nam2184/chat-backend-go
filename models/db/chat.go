package db

import "time"

const ChatTableName = "chats"

// Chat represents a conversation between two users
type Chat struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id"`                                        // Primary key
	User1ID         uint      `gorm:"not null;index;column:user1_id" db:"user1_id" json:"user1_id"`                                // ID of the first user
	User2ID         uint      `gorm:"not null;index;column:user2_id" db:"user2_id" json:"user2_id"`                                // ID of the second user
	Seen            bool      `gorm:"not null;column:seen" db:"seen" json:"seen"`                                                  // ID of the second user
	LastMessageTime time.Time `gorm:"index;column:last_message_time;default:null" db:"last_message_time" json:"last_message_time"` // Timestamp of the last message
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at" db:"created_at" json:"created_at"`                          // Auto-managed timestamp for creation
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at" db:"updated_at" json:"updated_at"`                          // Auto-managed timestamp for updates
}

func (m Chat) TableName() string {
	return ChatTableName
}

func (m Chat) Id() interface{} {
	return m.ID
}
