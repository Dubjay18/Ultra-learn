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
	_, err := r.db.Database("users").Collection("users").InsertOne(nil, user)
	if err != nil {
		return err
	}
	return nil

}

func (r *DefaultUserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Database("users").Collection("users").FindOne(nil, bson.D{
		{"email", email},
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.Database(database.DatabaseName).Collection("users").FindOne(nil, bson.D{
		{"_id", id},
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) UpdateUser(user *User) error {
	_, err := r.db.Database(database.DatabaseName).Collection("users").UpdateOne(nil, bson.D{
		{"_id", user.ID},
	}, bson.D{
		{"$set", bson.D{
			{"firstName", user.FirstName},
			{"lastName", user.LastName},
			{"email", user.Email},
			{"password", user.Password},
		}},
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *DefaultUserRepository) UpdateAvatar(id string, avatar string) error {
	_, err := r.db.Database(database.DatabaseName).Collection("users").UpdateOne(nil, bson.D{
		{"_id", id},
	}, bson.D{
		{"$set", bson.D{
			{"avatar", avatar},
		}},
	})

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
