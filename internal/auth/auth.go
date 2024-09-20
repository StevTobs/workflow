package aunt

import (
	"fmt"
	"net/http"
	"strings"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Secret key for signing JWT tokens (replace with your own secret)
var jwtSecret = []byte("123")

// Expected credentials
const (
	expectedUsername = "admin"
	expectedPassword = "secret"
)

var predefinedToken = ""
var tokenString = ""

func LoginHandler(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if loginRequest.Username != expectedUsername || loginRequest.Password != expectedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a token (this is a simplified example)
	predefinedToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoiU3VwZXIgQWRtaW4iLCJleHAiOjE"
	tokenString = predefinedToken
	// Set token in response cookies
	c.SetCookie("token", predefinedToken, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login succeeded"})
}

func LogoutHandler(c *gin.Context) {
	predefinedToken = ""
	c.SetCookie("token", "", -1, "/", "localhost", false, true) // Clear cookie
	c.JSON(http.StatusOK, gin.H{"message": "Logout succeeded"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Printf("Incoming token: %s\n", tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Check if the token is set and valid
		fmt.Printf("Expected token: %s\n", predefinedToken)
		if predefinedToken == "" || tokenString != predefinedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Login fails"})
			c.Abort()
			return
		}

		// Token is valid and the user is authenticated
		fmt.Println("Token is valid")
		c.Next()
	}
}
