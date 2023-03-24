package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	UserID   uint   `json:"usersid"`
}

func CreateProfile(db *gorm.DB, profile *Profile) (err error) {
	err = db.Create(&profile).Error

	if err != nil {
		return err
	}
	return nil
}

func GetProfiles(db *gorm.DB, profiles *[]Profile) (err error) {

	err = db.Find(&profiles).Error

	if err != nil {
		return err
	}
	return nil
}

func GetProfile(db *gorm.DB, profile *Profile, id string) (err error) {

	err = db.Where("id = ?", id).First(&profile).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProfileByUserID(db *gorm.DB, profile *Profile, userId uint) (err error) {

	err = db.Where("userid = ?", userId).First(&profile).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateProfile(db *gorm.DB, profile *Profile) (err error) {
	err = db.Save(&profile).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteProfile(db *gorm.DB, profile *Profile, id string) (err error) {
	err = db.Where("id = ?", id).Delete(&profile).Error
	if err != nil {
		return err
	}
	return nil
}
