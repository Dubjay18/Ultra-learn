package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/helper"
	"Ultra-learn/internal/repository"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) *errors.ApiError
	Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, *errors.ApiError)
}

type DefaultAuthService struct {
	repo *repository.DefaultUserRepository
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID string, role repository.Role) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"USER_ID": userID,
		"role":    role,
		"exp":     expirationTime,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *DefaultAuthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) *errors.ApiError {
	// Get the user data from the request
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		return &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
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
	err = a.repo.CreateUser(&repository.User{FirstName: user.FirstName,
		Email:    user.Email,
		Password: user.Password,
		LastName: user.LastName,
		Role:     "user",
		Avatar:   fmt.Sprintf("https://eu.ui-avatars.com/api/?name=%v+%v&size=250", user.FirstName, user.LastName),
		ID:       helper.GenerateUserId(),
	})
	if err != nil {
		return &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	return nil
}

func (a *DefaultAuthService) Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, *errors.ApiError) {
	// Get the user data from the request
	if err := c.ShouldBindJSON(&user); err != nil {
		return nil, &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
	}
	// Get the user from the database
	u, err := a.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	// Verify the user's password
	if !verifyPassword(user.Password, u.Password) {
		return nil, &errors.ApiError{
			Message:    errors.UnAuthorized,
			StatusCode: http.StatusUnauthorized,
			Error:      "Invalid email or password",
		}
	}
	// Generate a JWT token
	token, err := GenerateJWT(u.ID, u.Role)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	return &dto.LoginResponse{Token: token}, nil
}

func NewAuthService(repo *repository.DefaultUserRepository) AuthService {
	return &DefaultAuthService{
		repo: repo,
	}
}
