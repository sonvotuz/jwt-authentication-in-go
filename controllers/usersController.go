package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vnsonvo/jwt-authentication-in-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type body struct {
	Email    string
	Password string
}

func Signup(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		// Get the email/password from request body

		var requestBody = body{}
		if c.Bind(&requestBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read request body",
			})
			return
		}

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

func Login(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get the email and password from body request
		var requestBody = body{}

		if c.Bind(&requestBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read request body",
			})
			return
		}

		// look up requested user
		var user models.User
		db.First(&user, "email = ?", requestBody.Email)

		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// compare sent in pass with saved user pass hash
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// generate a jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		})

		// Sign and get the complete encoded token as a string using the secret
		secret := os.Getenv("SECRET")
		tokenString, err := token.SignedString([]byte(secret))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid to create token",
			})
			return
		}

		// set to cookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600, "", "", true, true)

		// send success response
		c.JSON(http.StatusOK, gin.H{})
	}
}

func Check(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not exists",
		})
		return
	}

	email := user.(models.User).Email
	c.JSON(http.StatusOK, gin.H{
		"message": map[string]string{"email": email},
	})
}
