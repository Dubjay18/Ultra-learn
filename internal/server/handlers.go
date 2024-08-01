package server

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"Ultra-learn/internal/helper"
	"Ultra-learn/internal/logger"
	errors2 "errors"
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

// RegisterUserHandler Register user handler
func (s *Server) RegisterUserHandler(c *gin.Context) {
	var user *dto.CreateUserRequest
	perr := helper.ParseRequestBody(c, &user)
	if perr != nil {
		return
	}
	resp, err := s.AuthService.CreateUser(c, user)
	if err != nil {
		// User creation failed
		c.JSON(err.StatusCode,
			err)
		return
	}
	// User created successfully
	Eerr := s.EmailService.SendSignUpEmail(user.Email, user.FirstName)
	if Eerr != nil {
		logger.Error("Error sending email: " + Eerr.Error())
	}

	helper.BuildSuccessResponse(c, http.StatusCreated, "User created successfully", resp)
}

// SignInUserHandler Sign in user handler
func (s *Server) SignInUserHandler(c *gin.Context) {
	if s.AuthService == nil {
		helper.BuildErrorResponse(c, http.StatusInternalServerError, "Internal server error", errors2.New("AuthService is not initialized"))
		return
	}
	var user dto.LoginRequest
	resp, err := s.AuthService.Login(c, &user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	helper.BuildSuccessResponse(c, http.StatusOK, "User logged in successfully", resp)
}

// user handlers

// GetUserDetailsHandler Get user details handler
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
	helper.BuildSuccessResponse(c, http.StatusOK, "User details retrieved successfully", user)
}

// UpdateUserDetailsHandler Update user details handler
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
			helper.BuildErrorResponse(c, http.StatusBadRequest, "Invalid request", errors2.New("Invalid request"))
			return
		}
	}

	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		helper.BuildErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userDetails, err := s.UserService.UpdateUserDetails(userID, &user)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	helper.BuildSuccessResponse(c, http.StatusOK, "User details updated successfully", userDetails)
}

func (s *Server) UpdateAvatarHandler(c *gin.Context) {
	userID := c.GetString("USER_ID")
	file, _, fileErr := c.Request.FormFile("avatar")
	if fileErr != nil {
		logger.Error("Error uploading file: " + fileErr.Error())
		helper.BuildErrorResponse(c, http.StatusBadRequest, "Error uploading file", fileErr)
		return
	}

	avatarUrl, err := s.UserService.UpdateAvatar(userID, file, c)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	helper.BuildSuccessResponse(c, http.StatusOK, "Avatar updated successfully", avatarUrl)
}
