package controllers

import (
	"net/http"
	"reg/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRoleRepo struct {
	DB *gorm.DB
}

func InitUserRoleController(db *gorm.DB) *UserRoleRepo {
	db.AutoMigrate(&models.User{}, &models.UserRole{})
	return &UserRoleRepo{DB: db}
}

func (repo UserRoleRepo) AddUserRole(c *gin.Context) {
	var userRole models.UserRole
	var userRoleApi models.UserRoleApi

	if c.Bind(&userRoleApi) == nil {
		var existingUser models.User

		models.GetUserByName(repo.DB, &existingUser, userRoleApi.Username)

		if existingUser.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "user is not found "})
			return
		}

		userRole.UserID = existingUser.ID
		userRole.Role = strings.ToUpper(userRoleApi.Role)

		err := models.AddUserRole(repo.DB, &userRole)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successfully added role to " + existingUser.Username})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

}

func (repo *UserRoleRepo) DeleteUserRole(c *gin.Context) {
	var userRole models.UserRole
	var userRoleApi models.UserRoleApi

	if c.Bind(&userRoleApi) == nil {
		var existingUser models.User
		models.GetUserByName(repo.DB, &existingUser, userRoleApi.Username)

		if existingUser.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		userRole.UserID = existingUser.ID
		userRole.Role = userRoleApi.Role

		err := models.DeleteUserRole(repo.DB, &userRole)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role Successfully deleted..."})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
}
