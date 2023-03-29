package main

import (
	"fmt"
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

	r1 := r.Group("/api/user")
	{
		r1.POST("/reg", userController.Register)
		r1.GET("/users", userController.GetUsers)
		r1.GET("/user/:username", userController.GetUser)
		r1.PUT("/changepassword", userController.ChangePassword)
		r1.PUT("/changeprofile", userController.ChangeProfile)
		r1.DELETE("/delete/:username", userController.DeleteUser)
	}

	r.Run()

}
