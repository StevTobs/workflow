package aunt

import (
	"fmt"
	"net/http"
	"strings"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Secret key for signing JWT tokens (replace with your own secret)
var jwtSecret = []byte("your_secret_key")

// Expected credentials
const (
	expectedUsername = "admin"
	expectedPassword = "secret"
)

var predefinedToken = ""

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
	predefinedToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoiU3VwZXIgQWRtaW4iLCJleHAiOjE"

	// Generate a token (this is a simplified example, consider using jwt-go library for JWT tokens)
	token := predefinedToken

	// Set token in response cookies
	c.SetCookie("token", "Bearer "+token, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login Succeed"})
}

func LogoutHandler(c *gin.Context) {
	predefinedToken = ""
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			fmt.Printf("tokenStrin : %s\n", tokenString)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		if predefinedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Login Fails"})
			c.Abort()
			return
		}

		// Token is valid and the user is authenticated
		fmt.Println("Token is valid")
		c.Next()
	}

}
