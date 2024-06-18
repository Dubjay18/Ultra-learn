package server

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) protectedHandler(c *gin.Context) {
	UserID := c.GetString("USER_ID")
	c.JSON(http.StatusOK, gin.H{"user_id": UserID})
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// auth handlers
// ...
// Register user handler
func (s *Server) registerUserHandler(c *gin.Context) {
	var user dto.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.ApiError{
			Message:    errors.ValidationError,
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return

	}
	err := s.authService.CreateUser(c, &user)
	if err != nil {
		// User creation failed
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}
