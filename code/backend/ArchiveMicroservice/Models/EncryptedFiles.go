package Models

import (
	"archivemicroservice/Config"
	"fmt"

	_ "gorm.io/driver/postgres"
)

func GetAllFiles(ef *[]EncryptedFile) (err error) {
	
	if err = Config.DB.Find(ef).Error; err != nil {
		return err
	}
	return nil
}

func CreateFile(ef *EncryptedFile) (err error) {
	if err = Config.DB.Create(ef).Error; err != nil {
		return err
	}
	return nil
}

func GetFile(ef *EncryptedFile, id int) (err error) {
	if err := Config.DB.Where("ID = ?", id).First(ef).Error; err != nil {
		return err
	}
	return nil
}

func GetFileByFilename(ef *EncryptedFile, filename string) (err error) {
	if err := Config.DB.Where("FILENAME = ?", filename).First(ef).Error; err != nil {
		return err
	}
	return nil
}

func UpdateFile(ef *EncryptedFile, id uint) (err error) {
	fmt.Println(ef)
	Config.DB.Save(ef)
	return nil
}

func DeleteUser(ef *EncryptedFile, id uint) (err error) {
	Config.DB.Where("id = ?", id).Delete(ef)
	return nil
}