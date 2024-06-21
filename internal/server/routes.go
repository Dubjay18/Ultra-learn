package server

import (
	"Ultra-learn/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/protected", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}), s.protectedHandler)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", s.registerUserHandler)
		v1.POST("/auth/login", s.signInUserHandler)
	}

	return r
}
