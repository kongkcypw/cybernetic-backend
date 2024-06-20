package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateCookieWithConfig(name string, value string, duration time.Duration) fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		Secure:   false, // Change this to true in production
		HTTPOnly: true,
		SameSite: "Lax",
	}
	return cookie
}
