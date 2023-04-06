package auth

import (
	"log"
	"net/http"
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
				// c.JSON(code, gin.H{
				// 	"access token": token,
				// 	"expires in ":  expire.Format(time.RFC1123),
				// })
				// persistance ....
				username, _ := c.Get("CURRENT_USERNAME")

				var userToken models.UserToken

				userToken.Token = token
				userToken.Username = username.(string)

				models.DeleteTokenByUsername(db, &userToken)

				err := models.SetToken(db, &userToken)
				if err == nil {
					c.JSON(http.StatusOK, gin.H{
						"token":     token,
						"expire in": expire.Format(time.RFC1123),
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "error while processing jwt token  .....",
					})
				}

			},

			LogoutResponse: func(c *gin.Context, code int) {
				var userToken models.UserToken

				oldToken, _ := c.Get("CURRENT_JWT_TOKEN")
				userToken.Token = oldToken.(string)

				err := models.DeleteToken(db, &userToken)
				if err == nil {
					c.JSON(http.StatusOK, gin.H{"message": "successfully deleted the token ..."})
				}

			},

			RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
				// c.JSON(code, gin.H{
				// 	"access token": token,
				// 	"expires in ":  expire.Format(time.RFC1123),
				// })

				var userToken models.UserToken

				oldToken, _ := c.Get("CURRENT_JWT_TOKEN")

				userToken.Token = oldToken.(string)
				models.DeleteToken(db, &userToken)

				userToken.Token = token
				err := models.SetToken(db, &userToken)

				if err == nil {
					c.JSON(code, gin.H{
						"access_token": token,
						"expires_in":   expire.Format(time.RFC3339),
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "problem in token processing ./.. "})
				}

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
