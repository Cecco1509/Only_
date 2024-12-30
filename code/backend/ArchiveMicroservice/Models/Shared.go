package Models

import (
	"archivemicroservice/Config"
)

func CreateShared(shared *Shared) (err error) {
	if err = Config.DB.Create(shared).Error; err != nil {
		return err
	}
	return nil
}

func GetShared(shared *Shared, fileId uint, sharedByUserId uint, sharedWithUserId uint) (err error) {
	if err := Config.DB.Where("FILE_ID = ? ", fileId).Where("SHARED_WITH_USER_ID = ? AND SHARED_BY_USER_ID = ?", sharedWithUserId, sharedByUserId).First(shared).Error; err != nil {
		return err
	}
	return nil
}

func DeleteShared(shared *Shared) (err error) {
	if err = Config.DB.Delete(shared).Error; err != nil {
		return err
	}
	return nil
}

func GetUserSharedWithFiles(ef *[]EncryptedFile, sharedByUserId uint, sharedWithUserId uint) (err error) {
	if err := Config.DB.Model(&EncryptedFile{}).Where("\"SHARED\".SHARED_BY_USER_ID = ? AND \"SHARED\".SHARED_WITH_USER_ID = ?", sharedByUserId, sharedWithUserId).Joins("JOIN \"SHARED\" ON \"ENCRYPTED_FILE\".ID = \"SHARED\".file_id").Find(ef).Error; err != nil {
			return err
		}
	return nil
}

func GetUserSharedFiles(ef *[]EncryptedFile, sharedByUserId uint) (err error) {
	if err := Config.DB.Where("SHARED_BY_USER_ID = ? ", sharedByUserId).Joins("JOIN \"ENCRYPTED_FILE\" ON \"ENCRYPTED_FILE\".USER_ID = ?", sharedByUserId).Find(ef).Error; err != nil {
			return err
		}
	return nil
}

func GetSharedWithUsers(shared *[]Shared, fileId uint) (err error) {
	if err := Config.DB.Where("FILE_ID = ?", fileId).Find(shared).Error; err != nil {
		return err
	}
	return nil
}

func GetShareUsers(results *[]Shared, sharedWithUserId uint) (err error) {
	if err := Config.DB.Where("shared_with_user_id = ?", sharedWithUserId).Group("shared_by_user_id, shared_by_username").Select("shared_by_user_id", "shared_by_username").Find(results).Error; err != nil {
		return err
	}
	return nil
}


func DeleteShares(shared *[]Shared, fileId uint) (err error) {
	if err = Config.DB.Where("FILE_ID = ?", fileId).Delete(shared).Error; err != nil {
		return err
	}
	return nil
}