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

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) CreateUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5)", user.FirstName, user.LastName, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil

}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT first_name, last_name, email, password, role FROM users WHERE email = $1", email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
