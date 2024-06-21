package services

import "Ultra-learn/internal/repository"

type UserService interface {
	getUserDetails(id string) (repository.User, error)
}
