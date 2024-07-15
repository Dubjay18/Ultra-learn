package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	Db *mongo.Client
}

var (
	DatabaseName = os.Getenv("DB_DATABASE")
	password     = os.Getenv("DB_PASSWORD")
	username     = os.Getenv("DB_USERNAME")
	port         = os.Getenv("DB_PORT")
	host         = os.Getenv("DB_HOST")
	dbInstance   *Service
)

func New() *Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err)

	}
	// Reuse Connection
	dbInstance = &Service{
		Db: client,
	}
	return dbInstance
}

func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.Db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
