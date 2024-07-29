package repository

import (
	"Ultra-learn/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	UpdateUser(user *User) error
	UpdateAvatar(id string, avatar string) error
}
type DefaultUserRepository struct {
	db *mongo.Client
}

func (r *DefaultUserRepository) CreateUser(user *User) error {
	err := database.InsertOne(r.db, "users", user)
	if err != nil {
		return err
	}
	return nil

}

func (r *DefaultUserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := database.GetOne(r.db, "users", bson.D{
		{"email", email},
	}, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) GetUserByID(id string) (*User, error) {
	var user User
	err := database.GetOne(r.db, "users", bson.D{
		{"_id", id},
	}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) UpdateUser(user *User) error {
	err := database.UpdateOne(r.db, "users", bson.D{
		{"_id", user.ID},
	}, bson.D{
		{"$set", bson.D{
			{"first_name", user.FirstName},
			{"last_name", user.LastName},
			{"email", user.Email},
			{"password", user.Password},
		},
		}})

	if err != nil {
		return err
	}
	return nil
}

func (r *DefaultUserRepository) UpdateAvatar(id string, avatar string) error {
	err := database.UpdateOne(r.db, "users", bson.D{
		{"_id", id},
	}, bson.D{
		{"$set", bson.D{
			{"avatar", avatar},
		},
		}})

	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *mongo.Client) *DefaultUserRepository {
	return &DefaultUserRepository{
		db: db,
	}
}
