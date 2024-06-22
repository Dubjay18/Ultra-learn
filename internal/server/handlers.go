package server

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

// Update user details handler
func (s *Server) updateUserDetailsHandler(c *gin.Context) {
	userID := c.GetString("USER_ID")

	var user dto.UpdateUserRequest
	var jsonData map[string]interface{}

	if err := c.ShouldBindBodyWith(&jsonData, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errors.ApiError{
			Message:    "Invalid request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	for key := range jsonData {
		switch key {
		case "first_name", "last_name", "email", "avatar":
			continue
		default:
			c.JSON(http.StatusBadRequest, errors.ApiError{
				Message:    "Invalid request",
				Error:      "Unexpected field " + key,
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}

	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errors.ApiError{
			Message:    "Invalid request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	userDetails, err := s.userService.UpdateUserDetails(userID, &user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "User details updated successfully",
		Data:       userDetails,
		StatusCode: http.StatusOK,
	})
}
