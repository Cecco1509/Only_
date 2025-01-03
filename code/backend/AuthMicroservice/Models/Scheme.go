package Models

import (
	"time"
)

type User struct {
    ID uint `gorm:"primaryKey autoIncrement" json:"id"`
    REGISTRATION_DT time.Time `gorm:"default:now()" json:"registration_dt"`
    LAST_LOGIN_DT time.Time `gorm:"default:null" json:"last_login_dt"` 
    USERNAME string `json:"username" gorm:"unique"`
    PASSHASH string `json:"-"`
    EMAIL string `json:"email" gorm:"unique"`
    FIRSTNAME string `json:"firstname"`
    LASTNAME string `json:"lastname"`
  }

func (b *User) TableName() string {
	return "USER"
}