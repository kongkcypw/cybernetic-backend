package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	helper "example/backend/helpers"
)

func Authenticate(c *fiber.Ctx) error {
	// Get the token from the header
	token := c.Get("authToken")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization token is required"})
	}

	// Verify the token
	claims, err := helper.VerifyToken(token)
	if err != "" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Set the user ID in the context
	c.Locals("userId", claims.UserId)
	c.Locals("email", claims.Email)

	// Continue the request if the token is valid
	return c.Next()
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GCP_CLIENT_ID"),
	ClientSecret: os.Getenv("GCP_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GCP_OAUTH_REDIRECT_URL"),
	Scopes:       []string{"profile", "email", "openid"},
	Endpoint:     google.Endpoint,
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

	// Important
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

	c.Locals("email", userInfo.Email)

	return c.Next()
}
