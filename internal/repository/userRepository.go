package repository

import (
	"Ultra-learn/internal/database"
	"Ultra-learn/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateAvatar(id string, avatar string) error
}
type DefaultUserRepository struct {
	db *database.Service
}

func (r *DefaultUserRepository) CreateUser(user *models.User) error {
	err := database.InsertOne(r.db, "users", user)
	if err != nil {
		return err
	}
	return nil

}

func (r *DefaultUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.GetOne(r.db, "users", bson.D{
		{Key: "email", Value: email},
	}, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := database.GetOne(r.db, "users", bson.D{
		{Key: "_id", Value: id},
	}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) UpdateUser(user *models.User) error {
	err := database.UpdateOne(r.db, "users", bson.D{
		{Key: "_id", Value: user.ID},
	}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "first_name", Value: user.FirstName},
			{Key: "last_name", Value: user.LastName},
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.Password},
		},
		}})

	if err != nil {
		return err
	}
	return nil
}

func (r *DefaultUserRepository) UpdateAvatar(id string, avatar string) error {
	err := database.UpdateOne(r.db, "users", bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "avatar", Value: avatar},
		},
		}})

	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *database.Service) *DefaultUserRepository {
	return &DefaultUserRepository{
		db: db,
	}
}
