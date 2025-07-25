package db

const NAMEAuth = "auth"

type Auth struct {
    Username string `json:"username" gorm:"column:username;type:varchar(100);not null;unique"`
    Password string `json:"password" gorm:"column:password;type:varchar(255);not null"`
    IsKhanh  bool   `json:"iskhanh" gorm:"column:is_khanh;type:boolean;default:false"`
    IsMy     bool   `json:"ismy" gorm:"column:is_my;type:boolean;default:false"`
}


func (m Auth) TableName() string {
	return NAMEAuth
}

func (m Auth) Id() interface{} {
	return m.Username
}
