package helpers

import (
	"go-api/backend-api/config"
	"time"                            

	"github.com/golang-jwt/jwt/v5"
)

// get JWT secret key from environment variable, if not set use default
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func GenerateToken(username string) string {

	// Set token expiration time, for example 60 minutes from now
	expirationTime := time.Now().Add(60 * time.Minute)

	// Create claims with username and expiration time
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create token with claims and sign it with the secret key with HS256 algorithm
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	// Return the generated token
	return token
}
