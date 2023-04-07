package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	Role   string `json:"role"`
	UserID uint   `json:"user_id"`
}

func GetRoleByUserId(db *gorm.DB, users *[]UserRole, id uint) (err error) {
	err = db.Where("user_id = ?", id).Find(&users).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func DeleteUserRole(db *gorm.DB, userRole *UserRole) (err error) {

	err = db.Unscoped().Where("user_id = ? AND role = ?", userRole.UserID, userRole.Role).Delete(&userRole).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func AddUserRole(db *gorm.DB, userRole *UserRole) (err error) {

	var userRoles []UserRole

	// getting the roles
	err = GetRoleByUserId(db, &userRoles, userRole.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {

			return err
		}
	}
	for _, item := range userRoles {
		if item.Role == userRole.Role {
			return nil
		}
	}

	// ! create new roles ..
	err = db.Create(userRole).Error

	if err != nil {
		log.Println("testing role ", err)
		return err
	}

	return nil
}
