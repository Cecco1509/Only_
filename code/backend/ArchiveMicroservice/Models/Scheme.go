package Models

import (
	"time"

	"gorm.io/gorm"
)

type EncryptedFile struct {
	gorm.Model `json:"-"`
  ID uint `gorm:"primaryKey autoIncrement" json:"id"` 
  CREATION_DT time.Time `gorm:"default:now()" json:"creation_dt"`
  LAST_RETREIVED_DT time.Time `gorm:"default:null" json:"last_login_dt"`
  LAST_UPDATE_DT time.Time `gorm:"default:null" json:"last_update_dt"`
  FILENAME string `json:"filename" gorm:"unique"`
  SALT string `json:"salt"`
  IV string `json:"iv"`
  EXTENSION string `json:"extension"`
  USER_ID uint `json:"user_id"`
}

type User struct {
  ID uint `json:"id"`
  REGISTRATION_DT time.Time `json:"registration_dt"`
  LAST_LOGIN_DT time.Time `json:"last_login_dt"` 
  USERNAME string `json:"username"`
  EMAIL string `json:"email"`
  FIRSTNAME string `json:"firstname"`
  LASTNAME string `json:"lastname"`
}


func (b *EncryptedFile) TableName() string {
	return "ENCRYPTED_FILE"
}