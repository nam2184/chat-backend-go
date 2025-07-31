package db

import "time"

const NAMEUser = "users"

// User mapped from table <users>
type User struct {
	ID        int64     `db:"id" json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	FirstName string    `db:"first_name" json:"first_name" gorm:"column:first_name;type:varchar(255)"`
	Surname   string    `db:"surname" json:"surname" gorm:"column:surname;type:varchar(255)"`
	Username  string    `db:"username" json:"username" gorm:"column:username;type:varchar(255);unique"`
	Email     string    `db:"email" json:"email" gorm:"column:email;type:varchar(255);unique"`
	CreatedAt time.Time `db:"created_at" json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (m User) TableName() string {
	return NAMEUser
}

func (m User) Id() interface{} {
	return m.ID
}
