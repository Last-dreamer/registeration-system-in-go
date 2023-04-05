package auth

import (
	"log"
	"os"
	"reg/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"
var identityName = "name"

func InitJwt(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	JWT_KEY := os.Getenv("JWT_KEY")
	JWT_REALM := os.Getenv("JWT_REALM")

	db.AutoMigrate(&models.User{}, &models.Profile{}, &models.UserToken{})

	authMiddleware, err := jwt.New(
		&jwt.GinJWTMiddleware{
			Realm:      JWT_REALM,
			Key:        []byte(JWT_KEY),
			Timeout:    7 * 24 * time.Hour,
			MaxRefresh: 10 * 24 * time.Hour,

			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*models.UserPayload); ok {
					return jwt.MapClaims{identityKey: v.Username, identityName: v.Fullname}
				}
				return jwt.MapClaims{}
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, gin.H{"code": code, "message": message})
			},

			Authorizator: func(data interface{}, c *gin.Context) bool {

				if _, ok := data.(*models.UserPayload); ok {
					return true
				}
				return false
			},

			Authenticator: func(c *gin.Context) (interface{}, error) {
				var loginVals models.Login

				if err := c.ShouldBind(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}

				username := loginVals.Username
				password := loginVals.Password

				var user models.User
				var profile models.Profile

				err := models.GetUserByName(db, &user, username)

				if err == nil {
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
					if err == nil {
						err := models.GetProfileByUserID(db, &profile, user.ID)

						if err == nil {
							c.Set("CURRENT_USERNAME", username)
							return &models.UserPayload{
								Username: username,
								Fullname: profile.FullName,
							}, nil
						}
					}
				}

				return nil, jwt.ErrFailedAuthentication

			},

			LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
				c.JSON(code, gin.H{
					"access token": token,
					"expires in ":  expire.Format(time.RFC1123),
				})
			},

			RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
				c.JSON(code, gin.H{
					"access token": token,
					"expires in ":  expire.Format(time.RFC1123),
				})
			},

			IdentityHandler: func(c *gin.Context) interface{} {
				claims := jwt.ExtractClaims(c)
				return &models.UserPayload{Username: claims[identityKey].(string),
					Fullname: claims[identityName].(string)}
			},
			IdentityKey:   identityKey,
			TokenLookup:   "header: Authorization, query: token, cookie: jw",
			TokenHeadName: "Bearer",
			TimeFunc:      time.Now,
		})

	if err != nil {
		log.Fatal("jwt error", err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("jwt init error", errInit.Error())
	}
	return authMiddleware, err
}
