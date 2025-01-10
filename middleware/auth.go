package middleware

import (
	"RemiAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token from cookies
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the auth_token cookie
		token, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not found"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for use in the handlers
		c.Set("user_id", claims["user_id"])

		// Proceed to the next handler
		c.Next()
	}
}
