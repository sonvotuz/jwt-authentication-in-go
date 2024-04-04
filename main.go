package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + port)
}
