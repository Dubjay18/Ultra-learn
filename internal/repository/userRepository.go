package repository

import (
	"database/sql"
)

type Role string

type User struct {
	ID        string `json:"id"db:"id"`
	Avatar    string `json:"avatar"db:"avatar"`
	FirstName string `json:"firstName"db:"first_name"`
	LastName  string `json:"lastName"db:"last_name"`
	Email     string `json:"email"db:"email"`
	Password  string `json:"password"db:"password"`
	Role      Role   `json:"role"db:"role"`
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}
type DefaultUserRepository struct {
	db *sql.DB
}

func (r *DefaultUserRepository) CreateUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, password, role,id,avatar) VALUES ($1, $2, $3, $4, $5,$6,$7)", user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.ID, user.Avatar)
	if err != nil {
		return err
	}
	return nil

}

func (r *DefaultUserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT first_name, last_name, email, password, role,id,avatar FROM users WHERE email = $1", email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.ID, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT first_name, last_name, email, password, role,id,avatar FROM users WHERE id = $1", id).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.ID, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) UpdateUser(user *User) error {
	_, err := r.db.Exec("UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, role = $5,avatar=$6 WHERE id = $7", user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.Avatar, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *sql.DB) *DefaultUserRepository {
	return &DefaultUserRepository{
		db: db,
	}
}
