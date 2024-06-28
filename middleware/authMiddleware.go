package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	helper "example/backend/helpers"
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
		c.Set("email", claims.Email)

		// Continue the request if the token is valid
		c.Next()
	}
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GCP_CLIENT_ID"),
	ClientSecret: os.Getenv("GCP_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GCP_OAUTH_REDIRECT_URL"),
	Scopes:       []string{"profile", "email", "openid"},
	Endpoint:     google.Endpoint,
}

func GoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("ClientID:", os.Getenv("GCP_CLIENT_ID"))
		fmt.Println("ClientSecret:", os.Getenv("GCP_CLIENT_SECRET"))
		fmt.Println("RedirectURL:", os.Getenv("GCP_OAUTH_REDIRECT_URL"))
		fmt.Println("Endpoint:", google.Endpoint)

		code := c.Query("code")
		fmt.Println("code:", code)
		if code == "" {
			c.JSON(400, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}

		// Important
		googleOauthConfig.ClientID = os.Getenv("GCP_CLIENT_ID")
		googleOauthConfig.ClientSecret = os.Getenv("GCP_CLIENT_SECRET")
		googleOauthConfig.RedirectURL = os.Getenv("GCP_OAUTH_REDIRECT_URL")

		token, err := googleOauthConfig.Exchange(c.Request.Context(), code)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to exchange token", "details": err.Error()})
			c.Abort()
			return
		}

		userInfo, err := helper.GetUserInfoFromGoogleOauthToken(token.AccessToken)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get user info"})
			c.Abort()
			return
		}

		// Set the email in the context
		c.Set("email", userInfo.Email)
		c.Next()
	}
}
