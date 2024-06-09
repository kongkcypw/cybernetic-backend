package controllers

import (

	// "log"
	// "strconv"
	// "net/http"
	// "time"
	// "fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	orm "example/backend/config"
	helper "example/backend/helpers"
	"example/backend/models"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)

		// validate user input
		validationError := helper.ValidateSignupInput(user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber)
		if validationError != "" {
			c.JSON(400, gin.H{"error": validationError})
			return
		}

		// Check if email already exists
		emailAlreadyExists := orm.DB().Where("email = ?", user.Email).First(&user).Error
		if emailAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Email already exists"})
			return
		}

		// Check if phone number already exists
		phoneAlreadyExists := orm.DB().Where("phoneNumber = ?", user.PhoneNumber).First(&user).Error
		if phoneAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Phone number already exists"})
			return
		}

		// Generate a unique user ID
		for {
			user.UserId = helper.GenerateUserId()
			if err := orm.DB().Where("userId = ?", user.UserId).First(&user).Error; err != nil {
				break // No user found with this ID, so it's unique
			}
		}

		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)

		// save the user
		dbErr := orm.DB().Create(&user).Error
		if dbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to signup"})
			return
		}
		c.JSON(200, user)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")

		if userId == "" {
			c.JSON(400, gin.H{"error": "userId is required"})
			return
		}

		var user models.User
		dbErr := orm.DB().Where("userId = ?", userId).First(&user).Error

		if dbErr != nil {
			c.JSON(404, gin.H{"error": "User not found!"})
			return
		}
		c.JSON(200, user)
	}
}
