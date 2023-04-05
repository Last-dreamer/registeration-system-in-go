package main

import (
	"fmt"
	"reg/auth"
	"reg/controllers"
	"reg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("something is wrong ")
	}
}

func main() {

	r := gin.Default()

	db := database.InitDb()

	userController := controllers.NewUserController(db)

	r1 := r.Group("/api")
	{
		r1.POST("/reg", userController.Register)
		r1.GET("/users", userController.GetUsers)
		r1.GET("/user/:username", userController.GetUser)
		r1.PUT("/changepassword", userController.ChangePassword)
		r1.PUT("/changeprofile", userController.ChangeProfile)
		r1.DELETE("/delete/:username", userController.DeleteUser)
	}

	basicAuth := auth.InitBasicAuth(db)
	r2 := r.Group("/user", basicAuth.BasicAuth())
	{
		r2.GET("/profile/:username", userController.GetUser)
	}

	jwtMiddleware, _ := auth.InitJwt(db)

	r.POST("/login", jwtMiddleware.LoginHandler)

	j := r.Group("/member", jwtMiddleware.MiddlewareFunc())
	{
		j.GET("/profile/:username", userController.GetUser)
	}

	r.Run()

}
