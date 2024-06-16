package middleware

import (
	helper "example/backend/helpers"

	"github.com/gofiber/fiber/v2"
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
	c.Locals("firstName", claims.FirstName)
	c.Locals("lastName", claims.LastName)
	c.Locals("email", claims.Email)
	c.Locals("phoneNumber", claims.PhoneNumber)

	// Continue the request if the token is valid
	return c.Next()
}
