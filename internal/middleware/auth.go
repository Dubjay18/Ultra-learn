package middleware

import (
	"Ultra-learn/internal/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
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

func AuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{Message: "Unauthorized",
			Error:      "Token is required",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ApiError{Message: "Internal server error",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{Message: "Unauthorized",
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{Message: "Token has expired",
					Error:      err.Error(),
					StatusCode: http.StatusUnauthorized,
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.ApiError{
				Message:    "Invalid token",
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		c.Set("USER_ID", claims["USER_ID"])
		c.Set("role", claims["role"])

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{Message: "Unauthorized",
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	//if !isLoggedIn(c) {
	//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
	//	return
	//}
	// Call the next handler
	c.Next()
}

func AccessControlMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the required role
		role := getUserRole(c)
		if !isRoleAllowed(role, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, errors.ApiError{Message: "Forbidden",
				Error:      "You do not have permission to access this resource",
				StatusCode: http.StatusForbidden,
			})
			return
		}
		// Call the next handler
		c.Next()
	}
}
