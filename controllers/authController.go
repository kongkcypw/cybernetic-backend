package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	database "example/backend/database"
	helper "example/backend/helpers"
	"example/backend/models"
)

func Signup(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	// Validate user input
	validationError := helper.ValidateSignupInput(user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber)
	if validationError != "" {
		return c.Status(400).JSON(fiber.Map{"error": validationError})
	}

	// Check if email already exists
	emailAlreadyExists := database.MysqlDB().Where("email = ?", user.Email).First(&user).Error
	if emailAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	// Check if phone number already exists
	phoneAlreadyExists := database.MysqlDB().Where("phoneNumber = ?", user.PhoneNumber).First(&user).Error
	if phoneAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Phone number already exists"})
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
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set the JWT token in a cookie
	tokenExpired, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRED"))
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    authToken,
		Expires:  time.Now().Add(time.Duration(tokenExpired) * time.Second),
		Path:     "/",
		Domain:   os.Getenv("SERVER_ENV"),
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"userId": user.UserId})
}

func Login(c *fiber.Ctx) error {
	// Read the JSON body into a user struct
	var userLogin models.UserLogin
	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate user input
	validationError := helper.ValidateLoginInput(userLogin.Email, userLogin.Password)
	if validationError != "" {
		return c.Status(400).JSON(fiber.Map{"error": validationError})
	}

	// Check if email exists
	var user models.User
	emailExists := database.MysqlDB().Where("email = ?", userLogin.Email).First(&user).Error
	if emailExists != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email not found"})
	}
	fmt.Println(user.Password)

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	authToken, _, generateTokenErr := helper.GenerateJWT(user.UserId, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
	if generateTokenErr != nil {
		return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	}

	// Set the JWT token in a cookie
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRED"))
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    authToken,
		Expires:  time.Now().Add(time.Duration(authTokenExpired) * time.Second),
		Path:     "/",
		Domain:   os.Getenv("SERVER_ENV"),
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"userId": user.UserId})
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
