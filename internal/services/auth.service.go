package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) *errors.ApiError
}

type DefaultAuthService struct {
	repo *repository.UserRepository
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

func (a *DefaultAuthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) *errors.ApiError {
	// Get the user data from the request
	if err := c.ShouldBindJSON(&user); err != nil {
		return &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
		}
	}
	// Hash the user's password
	hash, err := hashPassword(user.Password)
	if err != nil {
		return &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	user.Password = hash
	// Save the user to the database
	err = a.repo.CreateUser(user)
	if err != nil {
		return &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	return nil
}

func NewAuthService(repo *repository.UserRepository) AuthService {
	return &DefaultAuthService{
		repo: repo,
	}
}
