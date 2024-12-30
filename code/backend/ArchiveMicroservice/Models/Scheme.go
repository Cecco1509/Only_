package Models

import (
	"time"

	"gorm.io/gorm"
)

type EncryptedFile struct {
  ID uint `gorm:"primaryKey autoIncrement" json:"id"` 
  CREATED_AT time.Time `gorm:"default:now()" json:"creation_dt"`
  LAST_RETREIVED_DT time.Time `gorm:"default:null" json:"last_retrieved_dt"`
  UPDATE_AT time.Time `gorm:"default:null" json:"update_dt"`
  FILENAME string `json:"filename" gorm:"unique"`
  EXTENSION string `json:"extension"`
  USER_ID uint `json:"user_id"`
}

type Shared struct {
  gorm.Model `json:"-"`
  FILE_ID uint `json:"file_id"`
  SHARED_BY_USER_ID uint `json:"shared_by_user_id"`
  SHARED_BY_USERNAME string `json:"shared_by_username"`
  SHARED_WITH_USERNAME string `json:"shared_with_username"`
  SHARED_WITH_USER_ID uint `json:"shared_with_user_id"`
}


func (b *EncryptedFile) TableName() string {
	return "ENCRYPTED_FILE"
}

func (b *Shared) TableName() string {
  return "SHARED"
}