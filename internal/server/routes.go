package server

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/protected", services.AuthMiddleware(), services.AccessControlMiddleware([]string{"admin", "user"}), s.protectedHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) protectedHandler(c *gin.Context) {
	UserID := c.GetString("USER_ID")
	c.JSON(http.StatusOK, gin.H{"user_id": UserID})
}

func (s *Server) registerUserHandler(c *gin.Context) {
	var user dto.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.authService.CreateUser(c, &user)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
