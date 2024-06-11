package middleware

import (
	helper "example/backend/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the header
		token := c.Request.Header.Get("authToken")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Verify the token
		claims, err := helper.VerifyToken(token)
		if err != "" {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userId", claims.UserId)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("email", claims.Email)
		c.Set("phoneNumber", claims.PhoneNumber)

		// Continue the request if the token is valid
		c.Next()
	}
}
