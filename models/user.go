package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string
	Profile  Profile
	UserRole []UserRole
}

// type APIUser struct {
// 	Username string
// }

func CreateUser(db *gorm.DB, user User) (err error) {

	err = db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func GetUsers(db *gorm.DB, user *[]User) (err error) {
	err = db.Model(&User{}).Find(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByName(db *gorm.DB, user *User, username string) (err error) {

	// log.Println("user name ", username, ":::", user)
	err = db.Where("username = ?", username).Find(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func UpdateUser(db *gorm.DB, user *User) (err error) {

	//! todo: to check this specially ....
	// db.First(&user, id)
	// var updateUser User
	// db.Model(&user).Updates(User{username: "", passwor: "asdf"})

	err = db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *gorm.DB, user *User, username string) (err error) {

	GetUserByName(db, user, username)

	tx := db.Begin()

	if tx.Where("username = ? ", username).Delete(&User{}); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	if tx.Where("user_id = ? ", user.ID).Delete(&Profile{}); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	return tx.Commit().Error

}
