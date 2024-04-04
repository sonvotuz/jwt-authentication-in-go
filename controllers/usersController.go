package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vnsonvo/jwt-authentication-in-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		// Get the email/password from request body
		type body struct {
			Email    string
			Password string
		}

		var requestBody = body{}
		if c.Bind(&requestBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read request body",
			})
			return
		}

		fmt.Println(requestBody)
		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to hash password",
			})
			return
		}

		// create user
		user := models.User{Email: requestBody.Email, Password: string(hash)}

		result := db.Create(&user)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		// response
		c.JSON(http.StatusOK, gin.H{})
	}
}
