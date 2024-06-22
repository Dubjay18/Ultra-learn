package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/helper"
	"Ultra-learn/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserService interface {
	GetUserDetails(id string) (*dto.UserDetailsResponse, *errors.ApiError)
	UpdateUserDetails(id string, user *dto.UpdateUserRequest) (*dto.UserDetailsResponse, *errors.ApiError)
	UpdateAvatar(id string, file any, c *gin.Context) (string, *errors.ApiError)
}

type DefaultUserService struct {
	repo *repository.DefaultUserRepository
}

func (s *DefaultUserService) GetUserDetails(id string) (*dto.UserDetailsResponse, *errors.ApiError) {
	if id == "" {
		return nil, &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      "User ID is required",
		}

	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.UserNotFound,
			StatusCode: http.StatusNotFound,
			Error:      err.Error(),
		}
	}
	return &dto.UserDetailsResponse{
		ID:        user.ID,
		Avatar:    user.Avatar,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
	}, nil
}

func (s *DefaultUserService) UpdateUserDetails(id string, user *dto.UpdateUserRequest) (*dto.UserDetailsResponse, *errors.ApiError) {
	if id == "" {
		return nil, &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      "User ID is required",
		}
	}
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.UserNotFound,
			StatusCode: http.StatusNotFound,
			Error:      err.Error(),
		}
	}
	if user.FirstName != "" {
		u.FirstName = user.FirstName
	}
	if user.LastName != "" {
		u.LastName = user.LastName
	}

	err = s.repo.UpdateUser(u)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	return &dto.UserDetailsResponse{
		ID:        u.ID,
		Avatar:    u.Avatar,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      string(u.Role),
	}, nil
}
func (s *DefaultUserService) UpdateAvatar(id string, file any, c *gin.Context) (string, *errors.ApiError) {
	if id == "" {
		return "", &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      "User ID is required",
		}
	}
	if file == nil {
		return "", &errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      "Image is required",
		}
	}
	result, err := helper.UploadImage(file, id, c)
	if err != nil {
		return "", &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	err = s.repo.UpdateAvatar(id, result)
	if err != nil {
		return "", &errors.ApiError{
			Message:    errors.InternalServerError,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	return result, nil
}

func NewUserService(repo *repository.DefaultUserRepository) UserService {
	return &DefaultUserService{
		repo: repo,
	}
}
