package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vnsonvo/jwt-authentication-in-go/controllers"
	"github.com/vnsonvo/jwt-authentication-in-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDb() {
	dsn := os.Getenv("DB")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
}

func syncDatabase() {
	DB.AutoMigrate(&models.User{})
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectToDb()
	syncDatabase()
}

func main() {
	r := gin.Default()
	port := os.Getenv("PORT")

	r.POST("/signup", controllers.Signup(DB))

	r.Run(":" + port)
}
