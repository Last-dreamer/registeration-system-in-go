package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = ConnectDb()
	return Db
}

func ConnectDb() *gorm.DB {
	var err error
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database : error", err)
		return nil
	}
	return db
}
