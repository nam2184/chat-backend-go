package db

const EvaluationKeyTableName = "chats"

// Chat represents a conversation between two users
type EvaluationKeys struct {
	ChatID int64  `gorm:"primaryKey;column:chat_id" db:"chat_id" json:"chat_id"` // Primary key
	Key    []byte `gorm:"type:bytea;column:key" db:"key" json:"key"`
}

func (m EvaluationKeys) TableName() string {
	return EncryptedMessageTableName
}

func (m EvaluationKeys) Id() interface{} {
	return m.ChatID
}
