package services

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/repository"
	"net/http"
)

type UserService interface {
	GetUserDetails(id string) (*dto.UserDetailsResponse, *errors.ApiError)
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

func NewUserService(repo *repository.DefaultUserRepository) UserService {
	return &DefaultUserService{
		repo: repo,
	}
}
