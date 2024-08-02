package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/helper"
	"Ultra-learn/internal/models"
	"Ultra-learn/internal/repository"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// TODO: set up social logins
// TODO: set up email verification
// TODO: set up password reset
// TODO: set up user roles
// TODO: set up user permissions
type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.UserDetailsResponse, *errors.ApiError)
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

func GenerateJWT(userID string, role int) (string, error) {
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

func (a *DefaultAuthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.UserDetailsResponse, *errors.ApiError) {
	// Get the user data from the request
	//if err := c.ShouldBindJSON(&user); err != nil {
	//
	//	return &errors.ApiError{
	//		Message:    errors.ValidationError,
	//		StatusCode: http.StatusBadRequest,
	//		Error:      err.Error(),
	//	}
	//}
	// Hash the user's password
	//check if user already exists
	_, err := a.repo.GetUserByEmail(user.Email)
	if err == nil {
		return nil, &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      "User already exists",
		}
	}
	hash, err := hashPassword(user.Password)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed to Create User",
		}
	}
	user.Password = hash

	userReq := &models.User{FirstName: user.FirstName,
		Email:    user.Email,
		Password: user.Password,
		LastName: user.LastName,
		Role:     dto.RoleUser,
		Avatar:   fmt.Sprintf("https://eu.ui-avatars.com/api/?name=%v+%v&size=250", user.FirstName, user.LastName),
		ID:       helper.GenerateUserId(),
	}
	// Save the user to the database
	dbErr := a.repo.CreateUser(userReq)
	if dbErr != nil {
		return nil, &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      dbErr.Error(),
		}
	}
	return &dto.UserDetailsResponse{
		ID:        userReq.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    userReq.Avatar,
	}, nil
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
