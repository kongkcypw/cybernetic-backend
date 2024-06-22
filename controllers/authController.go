package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	config "example/backend/config"
	database "example/backend/database"
	helper "example/backend/helpers"
	"example/backend/models"
)

func Signup(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	user.Provider = "custom"

	// Validate user input
	validationError := helper.ValidateSignupInput(user.Email, user.Username, user.Password)
	if validationError != "" {
		return c.Status(400).JSON(fiber.Map{"error": validationError})
	}

	// Check if email already exists
	emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
	if emailAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists", "conflict": "email"})
	}

	// Check if username already exists
	usernameAlreadyExists := database.MysqlDB().Where("username = ?", user.Username).First(&user).Error
	if usernameAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists", "conflict": "username"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	// Save the user
	dbErr := database.MysqlDB().Create(&user).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to signup"})
	}

	// Generate JWT token
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set token expiration
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
	authTokenDuration := time.Duration(authTokenExpired*24) * time.Hour
	// Config and set cookie
	cookie := config.CreateCookieWithConfig("authToken", authToken, authTokenDuration)
	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{"userId": user.UserId, "email": user.Email})
}

func Login(c *fiber.Ctx) error {
	// Read the JSON body into a user struct
	var userLogin models.UserLogin
	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate user input
	validationError := helper.ValidateLoginInput(userLogin.Username, userLogin.Password)
	if validationError != "" {
		return c.Status(400).JSON(fiber.Map{"error": validationError})
	}

	// Check if email exists
	var user models.User
	usernameExists := database.MysqlDB().Where("username = ?", userLogin.Username).First(&user).Error
	if usernameExists != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username not found"})
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set token expiration
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
	authTokenDuration := time.Duration(authTokenExpired*24) * time.Hour
	// Config and set cookie
	cookie := config.CreateCookieWithConfig("authToken", authToken, authTokenDuration)
	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{"userId": user.UserId, "email": user.Email})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Second),
		Path:     "/",
		Domain:   os.Getenv("SERVER_ENV"),
		Secure:   false,
		HTTPOnly: true,
	})
	return c.Status(200).JSON(fiber.Map{"message": "Logged out successfully"})
}

func SignupWithGoogle(c *fiber.Ctx) error {
	var user models.User
	user.Email = c.Locals("email").(string)
	user.Provider = "google"

	// Check if email already exists
	emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
	if emailAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to signup"})
	}

	// Generate JWT token
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set token expiration
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
	authTokenDuration := time.Duration(authTokenExpired*24) * time.Hour
	// Config and set cookie
	cookie := config.CreateCookieWithConfig("authToken", authToken, authTokenDuration)
	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{"userId": user.UserId, "email": user.Email})
}

func LoginWithGoogle(c *fiber.Ctx) error {
	var user models.User
	user.Email = c.Locals("email").(string)
	user.Provider = "google"

	// Check if email exists
	emailExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
	if emailExists != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email not found"})
	}

	// Generate JWT token
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.Email)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set token expiration
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
	authTokenDuration := time.Duration(authTokenExpired*24) * time.Hour
	// Config and set cookie
	cookie := config.CreateCookieWithConfig("authToken", authToken, authTokenDuration)
	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{"userId": user.UserId, "email": user.Email})
}
