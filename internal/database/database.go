package database

import (
	"Ultra-learn/config"
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	Db *mongo.Database
}

var (
	DatabaseName = config.DatabaseName
	password     = config.DatabasePassword
	username     = config.DatabaseUsername
	port         = config.DatabasePort
	host         = config.DatabaseHost
	dbInstance   *Service
)

func New(dbName string, dbHost string, dbPort string) *Service {
	log.Println("Connecting to database")
	connectionString := fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	log.Println(connectionString)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)

	}
	// Reuse Connection
	dbInstance = &Service{
		Db: client.Database(dbName),
	}
	return dbInstance
}

func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.Db.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
