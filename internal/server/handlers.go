package server

import (
	"Ultra-learn/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) protectedHandler(c *gin.Context) {
	UserID, _ := c.Get("USER_ID")
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "Protected route",
		Data:       gin.H{"USER_ID": UserID},
		StatusCode: http.StatusOK,
	})
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "Hello World",
		Data:       resp,
		StatusCode: http.StatusOK,
	})
}

// auth handlers
// ...
// Register user handler
func (s *Server) registerUserHandler(c *gin.Context) {
	var user dto.CreateUserRequest
	err := s.authService.CreateUser(c, &user)
	if err != nil {
		// User creation failed
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusCreated,
		dto.ApiSuccessResponse{
			Message:    "User created successfully",
			StatusCode: http.StatusCreated,
		})

}

// Sign in user handler
func (s *Server) signInUserHandler(c *gin.Context) {
	var user dto.LoginRequest
	resp, err := s.authService.Login(c, &user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "User logged in successfully",
		Data:       resp,
		StatusCode: http.StatusOK,
	})
}

// user handlers
// ...
// Get user details handler
func (s *Server) getUserDetailsHandler(c *gin.Context) {
	userID := c.GetString("USER_ID")
	user, err := s.userService.GetUserDetails(userID)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "User details retrieved successfully",
		Data:       user,
		StatusCode: http.StatusOK,
	})
}
