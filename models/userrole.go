package models

import (
	"errors"

	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	Role   string `json:"role"`
	UserID uint   `json:"userid"`
}

func GetRoleByUserId(db *gorm.DB, user *[]UserRole, id uint) (err error) {
	err = db.Where("userid = ?", id).Find(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func DeleteUserRole(db *gorm.DB, userRole *UserRole) (err error) {

	err = db.Unscoped().Where("userid = ? AND role = ?", userRole.UserID, userRole.Role).Delete(&userRole).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func AddUserRole(db *gorm.DB, userRole *UserRole) (err error) {

	var userRoles []UserRole

	// getting the roles
	err = GetRoleByUserId(db, &userRoles, userRole.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for _, item := range userRoles {
		if item.Role == userRole.Role {
			return nil
		}
	}

	// ! create new roles ..
	err = db.Create(&userRole).Error

	if err != nil {
		return err
	}

	return nil
}
