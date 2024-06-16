package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest)
}

type DefaultauthService struct {
	repo *repository.UserRepository
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("USER_ID", claims["USER_ID"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if !isLoggedIn(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		// Call the next handler
		c.Next()
	}
}

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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verrifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"USER_ID": userID,
		"role":    role,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *DefaultauthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) {
	// Get the user data from the request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Hash the user's password
	hash, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	user.Password = hash
	// Save the user to the database
	err = a.repo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func NewAuthService(repo *repository.UserRepository) AuthService {
	return &DefaultauthService{
		repo: repo,
	}
}
