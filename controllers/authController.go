package controllers

import (
	"context"
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
		var user models.UserAuth
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

		// Save the user to sql database
		dbErr := database.MysqlDB().Create(&user).Error
		if dbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to signup"})
			return
		}

		// Save the user character to mongo database
		var character models.UserCharacter
		character.UserId = user.UserId
		character.CharacterName = ""
		character.HeighestLevel = 1
		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}
		collection := database.MongoDBOpenCollection(client, "cybernetic", "user_character")
		_, mongodbErr := collection.InsertOne(context.Background(), character)
		if mongodbErr != nil {
			c.JSON(500, gin.H{"error": "Failed to create character", "details": mongodbErr.Error()})
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
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_DOMAIN_FOR_COOKIE"), false, true)

		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email, "characterName": character.CharacterName})
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
		var user models.UserAuth
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
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_DOMAIN_FOR_COOKIE"), false, true)
		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the JWT token in the cookie
		// c.SetCookie("authToken", "", -1, "/", os.Getenv("SERVER_ENV"), false, true)
		c.SetCookie("authToken", "", -1, "/", os.Getenv("SERVER_DOMAIN_FOR_COOKIE"), false, true)
		c.JSON(200, gin.H{"message": "Logged out successfully"})

	}
}

func LoginWithGoogle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserAuth
		user.Email = c.GetString("email")

		// Check if email exists
		emailExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
		if emailExists == nil {
			if user.Provider != "google" {
				c.JSON(400, gin.H{"error": "You have already signed up with this email using another provider. Please login with your email and password."})
				return
			}
		}
		if emailExists != nil {
			user.Provider = "google"
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
			// Save the user character to mongo database
			var character models.UserCharacter
			character.UserId = user.UserId
			character.CharacterName = ""
			character.HeighestLevel = 1
			client := database.MongoDB()
			if client == nil {
				c.JSON(500, gin.H{"error": "Database connection is not initialized"})
				return
			}
			collection := database.MongoDBOpenCollection(client, "cybernetic", "user_character")
			_, mongodbErr := collection.InsertOne(context.Background(), character)
			if mongodbErr != nil {
				c.JSON(500, gin.H{"error": "Failed to create character", "details": mongodbErr.Error()})
				return
			}
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
		c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_DOMAIN_FOR_COOKIE"), false, true)
		c.JSON(200, gin.H{"userId": user.UserId, "email": user.Email})
	}
}

// net/http: invalid Cookie.Domain "http://localhost:3000"; dropping domain attribute
// SetCookie: The domain is set to an empty string to ensure it works on localhost.
// c.SetCookie("authToken", authToken, authTokenExpired, "/", os.Getenv("SERVER_DOMAIN_FOR_COOKIE"), false, true)
