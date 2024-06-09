package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
)

func GenerateUserId() string {
	randomNumber := rand.Intn(1e9) // Generates a random number between 0 and 999999999
	return fmt.Sprintf("u%09d", randomNumber)
}

func ValidateSignupInput(firstName, lastName, email, password, phoneNumber string) string {

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
