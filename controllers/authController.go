package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	database "example/backend/database"
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
		emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
		if emailAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Email already exists"})
			return
		}

		// Check if phone number already exists
		phoneAlreadyExists := database.MysqlDB().Where("phoneNumber = ?", user.PhoneNumber).First(&user).Error
		if phoneAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Phone number already exists"})
			return
		}

		// Generate a unique user ID
		for {
			user.UserId = helper.GenerateUserId()
			if err := database.MysqlDB().Where("userId = ?", user.UserId).First(&user).Error; err != nil {
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
		dbErr := database.MysqlDB().Create(&user).Error
		if dbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to signup"})
			return
		}

		// generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// set the JWT token in a cookie
		_, noCookie := c.Cookie("authToken")
		if noCookie != nil {
			tokenExpired, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRED"))
			c.SetCookie("authToken", authToken, tokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)
		}

		c.JSON(200, user.UserId)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the JSON body into a user struct
		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Validate user input
		validationError := helper.ValidateLoginInput(userLogin.Email, userLogin.Password)
		if validationError != "" {
			c.JSON(400, gin.H{"error": validationError})
			return
		}

		// Check if email exists
		var user models.User
		emailExists := database.MysqlDB().Where("email = ?", userLogin.Email).First(&user).Error
		if emailExists != nil {
			c.JSON(400, gin.H{"error": "email not found"})
			return
		}
		fmt.Println(user.Password)

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// Set the JWT token in a cookie
		cookie, noCookie := c.Cookie("authToken")
		authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRED"))
		if cookie != "" {
			c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)
		}
		if noCookie != nil {
			c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)
		}

		c.JSON(200, user.UserId)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("authToken", "", -1, "/", os.Getenv("SERVER_ENV"), false, true)
		c.JSON(200, gin.H{"message": "Logged out successfully"})
	}
}
