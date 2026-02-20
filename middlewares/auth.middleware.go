package middlewares

import (
	"net/http"                      
	"go-api/backend-api/config" 
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// get JWT secret key from environment variable, if not set use default
// if doesn't work, you can hardcode it for testing purposes, but make sure to change it in production
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// take token from Authorization header
		tokenString := c.GetHeader("Authorization")

		// if token is empty, return unauthorized
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort() // Hold other requests until this one is processed
			return
		}

		// delete "Bearer " prefix if it exists, "Bearer <token>"
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Define a struct to hold the claims (payload) from the JWT token
		claims := &jwt.RegisteredClaims{}

		// Parse token and validate it with the secret key
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// return the secret key for validation
			return jwtKey, nil
		})

		// if there's an error and retrun token is invalid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort() // Hentikan request
			return
		}

		// Set username in context, you can access it in the handler with c.Get("username")
		c.Set("username", claims.Subject)

		// Continue to the next handler
		c.Next()
	}
}
