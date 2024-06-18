package server

import (
	"Ultra-learn/internal/repository"
	"Ultra-learn/internal/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"Ultra-learn/internal/database"
)

type Server struct {
	port        int
	db          *database.Service
	authService services.AuthService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbInstance := database.New()

	userRepo := repository.NewUserRepository(dbInstance.Db) // Pass the dbInstance to the UserRepository
	authService := services.NewAuthService(userRepo)        // Pass the UserRepository to the AuthService
	NewServer := &Server{
		port:        port,
		db:          dbInstance,
		authService: authService,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
