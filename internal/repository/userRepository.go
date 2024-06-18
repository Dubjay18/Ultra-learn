package repository

import (
	"Ultra-learn/internal/dto"
	"database/sql"
)

type Role string

type User struct {
	FirstName string `json:"firstName"db:"first_name"`
	LastName  string `json:"lastName"db:"last_name"`
	Email     string `json:"email"db:"email"`
	Password  string `json:"password"db:"password"`
	Role      Role   `json:"role"db:"role"`
}

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) CreateUser(user *dto.CreateUserRequest) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5)", user.FirstName, user.LastName, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil

}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
