package auth

import (
	"net/http"
	"reg/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HelperRepo struct {
	DB *gorm.DB
}

func InitHelper(db *gorm.DB) *HelperRepo {
	db.AutoMigrate(&models.User{})
	return &HelperRepo{DB: db}
}

func (repo *HelperRepo) VerifyToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	token := ParseToken(authHeader)

	if token != "" {
		var userToken models.UserToken
		userToken.Token = token

		err := models.GetToken(repo.DB, &userToken, userToken.Token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized ...."})
			return
		}
		c.Set("CURRENT_JWT_TOKEN", token)
	} else {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized 123...."})
		return
	}

}

func ParseToken(auth string) (token string) {
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)

		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}
