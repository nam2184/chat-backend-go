package db

const NAMEAuth = "auth"

type Auth struct {
	Username string `json:"username" gorm:"column:username;type:varchar(100);not null;unique"`
	Password string `json:"password" gorm:"column:password;type:varchar(255);not null"`
}

func (m Auth) TableName() string {
	return NAMEAuth
}

func (m Auth) Id() interface{} {
	return m.Username
}
