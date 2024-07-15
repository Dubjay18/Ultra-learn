package server

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.Db.Health())
}

func (s *Server) ProtectedHandler(c *gin.Context) {
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
func (s *Server) RegisterUserHandler(c *gin.Context) {
	if s.AuthService == nil {
		c.JSON(http.StatusInternalServerError, errors.ApiError{
			Message:    "Internal server error",
			Error:      "AuthService is not initialized",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	user := dto.CreateUserRequest{}
	err := s.AuthService.CreateUser(c, &user)
	if err != nil {
		// User creation failed
		c.JSON(
			err.StatusCode,
			err)
		return
	}
	c.JSON(http.StatusCreated,
		dto.ApiSuccessResponse{
			Message:    "User created successfully",
			StatusCode: http.StatusCreated,
		})

}

// Sign in user handler
func (s *Server) SignInUserHandler(c *gin.Context) {
	if s.AuthService == nil {
		c.JSON(http.StatusInternalServerError, errors.ApiError{
			Message:    "Internal server error",
			Error:      "AuthService is not initialized",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	var user dto.LoginRequest
	resp, err := s.AuthService.Login(c, &user)
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
func (s *Server) GetUserDetailsHandler(c *gin.Context) {
	if s.AuthService == nil {
		c.JSON(http.StatusInternalServerError, errors.ApiError{
			Message:    "Internal server error",
			Error:      "AuthService is not initialized",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	userID := c.GetString("USER_ID")
	user, err := s.UserService.GetUserDetails(userID)
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
func (s *Server) UpdateUserDetailsHandler(c *gin.Context) {
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

	userDetails, err := s.UserService.UpdateUserDetails(userID, &user)
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

func (s *Server) UpdateAvatarHandler(c *gin.Context) {
	userID := c.GetString("USER_ID")
	file, _, fileErr := c.Request.FormFile("avatar")
	if fileErr != nil {
		c.JSON(http.StatusBadRequest, errors.ApiError{
			Message:    "Invalid request",
			Error:      fileErr.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	avatarUrl, err := s.UserService.UpdateAvatar(userID, file, c)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Message:    "Avatar updated successfully",
		StatusCode: http.StatusOK,
		Data:       avatarUrl,
	})
}
