package server

import (
	"Ultra-learn/config"
	"Ultra-learn/internal/logger"
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
	Port         int
	Db           *database.Service
	AuthService  services.AuthService
	UserService  services.UserService
	EmailService services.EmailService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbInstance := database.New(config.DatabaseName, config.DatabaseHost, config.DatabasePort)
	logger.Init()
	userRepo := repository.NewUserRepository(dbInstance) // Pass the dbInstance to the UserRepository
	authService := services.NewAuthService(userRepo)     // Pass the UserRepository to the AuthService
	userService := services.NewUserService(userRepo)     // Pass the UserRepository to the UserService
	emailService := services.NewEmailService()           // Pass the UserRepository to the EmailService
	NewServer := &Server{
		Port:         port,
		Db:           dbInstance,
		AuthService:  authService,
		UserService:  userService,
		EmailService: emailService,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
