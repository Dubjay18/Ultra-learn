package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func isLoggedIn(c *gin.Context) bool {
	userID := c.GetString("USER_ID")
	return userID != ""
}

func getUserRole(c *gin.Context) string {
	// Get the user's role from the session or database
	// Example: get the role from the session
	role := c.GetString("role")
	return role
}

func isRoleAllowed(role string, allowedRoles []string) bool {
	// Check if the user's role is allowed
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			return
		}
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": "Internal Server Error"})
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("USER_ID", claims["USER_ID"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			return
		}
		if !isLoggedIn(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			return
		}
		// Call the next handler
		c.Next()
	}
}
func AccessControlMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the required role
		role := getUserRole(c)
		if !isRoleAllowed(role, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"errors": "Forbidden"})
			return
		}
		// Call the next handler
		c.Next()
	}
}
