package auth

import (
	"net/http"
	"reg/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BasicAuthRepo struct {
	Db *gorm.DB
}

func InitBasicAuth(db *gorm.DB) *BasicAuthRepo {
	db.AutoMigrate(&models.User{})
	return &BasicAuthRepo{Db: db}
}

func (repo BasicAuthRepo) BasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		var user models.User
		var err error

		username, password, hasAuth := c.Request.BasicAuth()

		if hasAuth {
			err = repo.Db.Where("username = ?", username).First(&user).Error
			if err == nil {
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == nil {
					c.Next()
					return
				}
			}
		}
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthurized ....."})

	}
}
