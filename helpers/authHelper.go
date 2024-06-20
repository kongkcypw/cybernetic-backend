package helpers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
)

func GenerateUserId() string {
	randomNumber := rand.Intn(1e9) // Generates a random number between 0 and 999999999
	return fmt.Sprintf("u%09d", randomNumber)
}

func ValidateSignupInput(firstName, lastName, email, username, password, phoneNumber string) string {

	// Validate first name
	if firstName == "" {
		return "First name is required"
	}

	// Validate last name
	if lastName == "" {
		return "Last name is required"
	}

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "Invalid email address"
	}

	if username == "" {
		return "Username is required"
	}

	// Validate password
	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	// Validate phone number
	phoneNumberRegex := regexp.MustCompile(`^\d{10}$`)
	if !phoneNumberRegex.MatchString(phoneNumber) {
		return "Invalid phone number"
	}

	return ""
}

func ValidateLoginInput(username, password string) string {
	// Validate email
	if username == "" {
		return "Username is required"
	}

	// Validate password
	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	return ""
}

type GoogleProfile struct {
	Email string `json:"email"`
}

func GetUserInfoFromGoogleOauthToken(accessToken string) (map[string]interface{}, error) {
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send HTTP request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response body
	// var userInfo map[string]interface{}
	var userInfo map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
