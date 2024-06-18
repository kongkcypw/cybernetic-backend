package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

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
	validationError := helper.ValidateSignupInput(user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.PhoneNumber)
	if validationError != "" {
		return c.Status(400).JSON(fiber.Map{"error": validationError})
	}

	// Check if username already exists
	usernameAlreadyExists := database.MysqlDB().Where("username = ?", user.Username).First(&user).Error
	if usernameAlreadyExists == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
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

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GCP_CLIENT_ID"),
	ClientSecret: os.Getenv("GCP_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GCP_OAUTH_REDIRECT_URL"),
	Scopes:       []string{"profile", "email", "openid"},
	Endpoint:     google.Endpoint,
}

func GoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("state")
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	fmt.Println("ClientID:", os.Getenv("GCP_CLIENT_ID"))
	fmt.Println("ClientSecret:", os.Getenv("GCP_CLIENT_SECRET"))
	fmt.Println("RedirectURL:", os.Getenv("GCP_OAUTH_REDIRECT_URL"))
	fmt.Println("Endpoint:", google.Endpoint)

	code := c.Query("code")
	fmt.Println("code:", code)
	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid code"})
	}

	googleOauthConfig.ClientID = os.Getenv("GCP_CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("GCP_CLIENT_SECRET")
	googleOauthConfig.RedirectURL = os.Getenv("GCP_OAUTH_REDIRECT_URL")

	token, err := googleOauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to exchange token", "details": err.Error()})
	}

	userInfo, err := helper.GetUserInfoFromGoogleOauthToken(token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user info"})
	}

	// authToken, _, generateTokenErr := helper.GenerateJWT(userInfo.given_name, userInfo.FirstName, userInfo.LastName, userInfo.Email, userInfo.PhoneNumber)
	// if generateTokenErr != nil {
	// 	return c.Status(500).JSON(fiber.Map{"error": generateTokenErr.Error()})
	// }

	return c.Status(200).JSON(fiber.Map{"profile": userInfo})
	// return c.Status(200).JSON(fiber.Map{"message": "Google callback"})
}
