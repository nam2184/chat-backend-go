package db

import "time"

const NAMEUser = "users"

// User mapped from table <users>
type User struct {
	ID        int64    `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	Surname   string    `db:"surname" json:"surname"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	IsMy      bool      `db:"is_my" json:"is_my"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsKhanh   bool      `db:"is_khanh" json:"is_khanh"`
}

func (m User) TableName() string {
	return NAMEUser
}

func (m User) Id() interface{} {
	return m.ID
}
