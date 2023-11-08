package utils

import (
	"log"
	"os"
	"strings"

	"github.com/jmcerez0/fiber-demo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CreateDB() {
	dsn := strings.Split(os.Getenv("DATABASE_URL"), "/")[0] + "/"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database server.")
	}

	db.Exec("CREATE DATABASE IF NOT EXISTS fiber_demo")

	DB, _ := db.DB()
	DB.Close()
}

func ConnectToDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error establishing a database connection.")
	}

	DB = db
}

func MigrateSchema() {
	DB.AutoMigrate(&models.User{})
}
