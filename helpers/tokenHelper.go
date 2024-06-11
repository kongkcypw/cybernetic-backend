package helpers

import (
	"log"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserId      string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	jwt.StandardClaims
}

func GenerateJWT(userId string, firstName string, lastName string, email string, phoneNumber string) (string, string, error) {

	// Set Expired time
	authTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_AUTH_TOKEN_EXPIRED"))
	refreshTokenExpired, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRED"))

	// Set the token claims
	claims := &TokenClaims{
		UserId:      userId,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(authTokenExpired),
		},
	}
	refreshClaims := &TokenClaims{
		UserId:      userId,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(refreshTokenExpired),
		},
	}

	// Generate JWT token
	token, tokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))
	if tokenErr != nil {
		log.Panic(tokenErr)
		return "", "", tokenErr
	}
	refreshToken, refreshErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))
	if refreshErr != nil {
		log.Panic(refreshErr)
		return "", "", refreshErr
	}

	return token, refreshToken, nil
}

func VerifyToken(tokenString string) (claims *TokenClaims, msg string) {
	// Parse the token
	token, tokenErr := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_TOKEN_SECRET")), nil
		},
	)

	if tokenErr != nil {
		msg = tokenErr.Error()
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		msg = "token is invalid"
		return
	}

	if claims.ExpiresAt < 0 {
		msg = "token is expired"
		return
	}

	return claims, ""
}
