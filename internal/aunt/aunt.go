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

// Predefined token for authentication (example, replace with dynamic tokens in production)

// func splitJWT(token string) (header, payload, signature string, err error) {
// 	// Split the token by dot (`.`)
// 	fmt.Printf("token :%s\n", token)
// 	parts := strings.Split(token, ".")
// 	if len(parts) != 3 {
// 		return "", "", "", fmt.Errorf("invalid token format")
// 	}

// 	return parts[0], parts[1], parts[2], nil
// }

// LoginHandler handles user login and returns a token if credentials are valid
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

		// Extract token from the header
		// tokenString = tokenString[7:] // Remove "Bearer " prefix

		// Check if the token has the correct number of segments
		// header, payload, signature, err := splitJWT(tokenString)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// } else {
		// 	fmt.Println("Header:", header)
		// 	fmt.Println("Payload:", payload)
		// 	fmt.Println("Signature:", signature)
		// }

		// Decode and verify the token
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	// Ensure the token method conforms to the expected signing method
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// 	}
		// 	return jwtSecret, nil
		// })

		// if err != nil {
		// 	// Print the detailed error message
		// 	fmt.Printf("Error during token parsing: %s\n", err)
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// if !token.Valid {
		// 	fmt.Println("Token is not valid")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }{
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
