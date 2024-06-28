package controllers

import (
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

		// Read the JSON body into a user struct
		var user models.User
		c.BindJSON(&user)

		user.Provider = "custom"

		// Validate user input
		validationError := helper.ValidateSignupInput(user.Email, user.Username, user.Password)
		if validationError != "" {
			c.JSON(400, gin.H{"error": validationError})
			return
		}

		// Check if email already exists
		emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
		if emailAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Email already exists", "conflict": "email"})
			return
		}

		// Check if username already exists
		usernameAlreadyExists := database.MysqlDB().Where("username = ?", user.Username).First(&user).Error
		if usernameAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Username already exists", "conflict": "username"})
			return
		}

		// Generate a unique user ID
		for {
			user.UserId = helper.GenerateUserId()
			if err := database.MysqlDB().Where("userId = ?", user.UserId).First(&user).Error; err != nil {
				break // No user found with this ID, so it's unique
			}
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)

		// Save the user
		dbErr := database.MysqlDB().Create(&user).Error
		if dbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to signup"})
			return
		}

		// Generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// Set token expiration
		authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)

		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the JSON body into a user struct
		var userLogin models.UserLogin
		if err := c.BindJSON(&userLogin); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Validate user input
		validationError := helper.ValidateLoginInput(userLogin.Username, userLogin.Password)
		if validationError != "" {
			c.JSON(400, gin.H{"error": validationError})
			return
		}

		// Check if email exists
		var user models.User
		usernameExists := database.MysqlDB().Where("username = ?", userLogin.Username).First(&user).Error
		if usernameExists != nil {
			c.JSON(400, gin.H{"error": "Username not found"})
			return
		}

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// Set token expiration
		authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
		// Config and set cookie
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)

		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the JWT token in the cookie
		c.SetCookie("authToken", "", -1, "/", os.Getenv("SERVER_ENV"), false, true)
		c.JSON(200, gin.H{"message": "Logged out successfully"})

	}
}

func SignupWithGoogle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		user.Email = c.GetString("email")
		user.Provider = "google"

		// Check if email already exists
		emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
		if emailAlreadyExists == nil {
			c.JSON(400, gin.H{"error": "Email already exists"})
			return

		}

		// Generate a unique user ID
		for {
			user.UserId = helper.GenerateUserId()
			if err := database.MysqlDB().Where("userId = ?", user.UserId).First(&user).Error; err != nil {
				break // No user found with this ID, so it's unique
			}
		}
		// Save the user
		dbErr := database.MysqlDB().Create(&user).Error
		if dbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to signup"})
			return
		}

		// Generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// Set token expiration
		authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
		// Config and set cookie
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)

		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}

func LoginWithGoogle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		user.Email = c.GetString("email")
		user.Provider = "google"

		// Check if email exists
		emailExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
		if emailExists != nil {
			c.JSON(400, gin.H{"error": "Email not found"})
			return
		}

		// Generate JWT token
		authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
		if generateTokenErr != nil {
			c.JSON(500, gin.H{"error": generateTokenErr.Error()})
			return
		}

		// Set token expiration
		authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
		// Config and set cookie
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_ENV"), false, true)
		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}
