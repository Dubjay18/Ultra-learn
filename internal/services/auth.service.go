package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/helper"
	"Ultra-learn/internal/models"
	"Ultra-learn/internal/repository"
	errors2 "errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// TODO: set up social logins
// TODO: set up email verification
// TODO: set up password reset
// TODO: set up user roles
// TODO: set up user permissions
type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.UserDetailsResponse, error)
	Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, error)
	CheckIFUserExists(c *gin.Context, email string) (bool, error)
	SocialLogin(c *gin.Context, user goth.User) (*dto.LoginResponse, error)
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

func (a *DefaultAuthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.UserDetailsResponse, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
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
		return nil, dbErr
	}
	return &dto.UserDetailsResponse{
		ID:        userReq.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    userReq.Avatar,
	}, nil
}

func (a *DefaultAuthService) Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get the user from the database
	u, err := a.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	// Verify the user's password
	if !verifyPassword(user.Password, u.Password) {
		return nil, errors2.New("invalid password or email")
	}
	// Generate a JWT token
	token, err := GenerateJWT(u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Token: token}, nil
}

func (a *DefaultAuthService) CheckIFUserExists(c *gin.Context, email string) (bool, error) {
	_, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return true, nil

}
func (a *DefaultAuthService) SocialLogin(c *gin.Context, user goth.User) (*dto.LoginResponse, error) {
	// Check if the user already exists by email
	exists, err := a.CheckIFUserExists(c, user.Email)
	if err != nil && !errors2.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	var u *models.User
	if !exists {
		// Create a new user if not found
		u = &models.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Avatar:    user.AvatarURL,
			Role:      dto.RoleUser,
			ID:        helper.GenerateUserId(),
			Provider:  user.Provider,
		}

		if err := a.repo.CreateUser(u); err != nil {
			return nil, err
		}
	} else {
		// Fetch the user if they already exist
		u, err = a.repo.GetUserByEmail(user.Email)
		if err != nil {
			return nil, err
		}
	}

	// Generate a JWT token for the user
	token, err := GenerateJWT(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{Token: token}, nil
}

func NewAuthService(repo *repository.DefaultUserRepository) AuthService {
	return &DefaultAuthService{
		repo: repo,
	}
}
