package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"reg/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) (user *UserRepo) {
	db.AutoMigrate(&models.User{}, &models.Profile{})
	return &UserRepo{DB: db}
}

func (repo *UserRepo) Register(c *gin.Context) {
	var user models.User
	var profile models.Profile
	var reg models.UserRegister

	if c.BindJSON(&reg) != nil {

		var existingUser models.User
		models.GetUserByName(repo.DB, &existingUser, reg.UserName)

		if existingUser.ID > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is already exist ..."})
			return
		}

		user.Username = reg.UserName
		profile.FullName = reg.FullName
		profile.Email = reg.Email

		encPass, _ := bcrypt.GenerateFromPassword([]byte(reg.Password), 10)
		user.Password = string(encPass)

		err2 := models.CreateUser(repo.DB, user)

		if err2 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successfully create a user ", "data": user})
	} else {
		fmt.Println(user)
		c.JSON(http.StatusBadRequest, user)
	}

}

func (repo *UserRepo) GetUsers(c *gin.Context) {
	var profiles []models.Profile

	err := models.GetProfiles(repo.DB, &profiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": profiles})
}

func (repo *UserRepo) GetUser(c *gin.Context) {
	username, _ := c.Params.Get("username")
	var user models.User

	err := models.GetUserByName(repo.DB, &user, username)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// c.JSON(http.StatusNotFound, gin.H{"messsage": "user not found "})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	var profile models.Profile

	err2 := models.GetProfileByUserID(repo.DB, &profile, uint(user.ID))
	if err2 != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"message": err2})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (repo *UserRepo) ChangeProfile(c *gin.Context) {
	var err error
	var user models.User
	var updatedPassword models.NewPassword

	if c.BindJSON(&updatedPassword) != nil {

		// ! check the user ..
		err = models.GetUserByName(repo.DB, &user, updatedPassword.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// ! verify the current password ...
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updatedPassword.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// ! new password
		passwordBytes, _ := bcrypt.GenerateFromPassword([]byte(updatedPassword.NewPassword), 10)
		user.Password = string(passwordBytes)

		err = models.UpdateUser(repo.DB, &user)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password successfully changed ..."})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid request"})
	}
}