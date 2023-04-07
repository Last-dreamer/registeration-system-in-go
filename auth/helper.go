package auth

import (
	"net/http"
	"reg/models"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
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

func (repo *HelperRepo) CheckRoles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, _ := c.Get("JWT_PAYLOAD")
		claims := payload.(jwt.MapClaims)

		var user models.User
		err := models.GetUserByName(repo.DB, &user, claims["id"].(string))

		if err == nil {
			var userRole []models.UserRole
			err := models.GetRoleByUserId(repo.DB, &userRole, user.ID)
			if err == nil {

				for _, userrole := range userRole {
					for _, role := range roles {
						if strings.EqualFold(userrole.Role, role) {
							c.Next()
							return
						}
					}
				}
			}
		}
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized ....."})

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
