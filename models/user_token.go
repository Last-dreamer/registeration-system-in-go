package models

import (
	"errors"

	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model
	Username string
	Token    string `gorm:"unique"`
}

func SetToken(db *gorm.DB, userToken *UserToken) (err error) {
	err = GetToken(db, userToken, userToken.Token)
	if err == nil {
		err = db.Create(&userToken).Error
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create(&userToken).Error
		}
	}

	return err
}

func GetToken(db *gorm.DB, userToken *UserToken, token string) (err error) {
	err = db.Where("token = ?", token).First(&userToken).Error
	return err
}

func DeleteToken(db *gorm.DB, userToken *UserToken) (err error) {

	err = db.Unscoped().Where("token = ?", userToken.Token).Delete(&UserToken{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil

	}
	return err
}

func DeleteTokenByUsername(db *gorm.DB, userToken *UserToken) (err error) {

	err = db.Unscoped().Where("username = ?", userToken.Username).Delete(&UserToken{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil

	}
	return err
}
