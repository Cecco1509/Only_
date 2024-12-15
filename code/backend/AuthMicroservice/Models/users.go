package Models

import (
	"authmicroservice/Config"
	"fmt"

	_ "gorm.io/driver/postgres"
)

func GetAllUsers(u *[]User) (err error) {
	
	if err = Config.DB.Find(u).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(u *User) (err error) {
	if err = Config.DB.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(u *User, username string) (err error) {
	if err := Config.DB.Where("USERNAME = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(u *User, id uint) (err error) {
	fmt.Println(u)
	Config.DB.Save(u)
	return nil
}

func DeleteUser(u *User, id string) (err error) {
	Config.DB.Where("id = ?", id).Delete(u)
	return nil
}