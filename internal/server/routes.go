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
	ar := v1.Group("/auth")
	ur := v1.Group("/user", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}))
	{

		ar.POST("/register", s.registerUserHandler)
		ar.POST("/login", s.signInUserHandler)
		ur.GET("/details", s.getUserDetailsHandler)
		ur.PUT("/details", s.updateUserDetailsHandler)
		ur.POST("/avatar", s.updateAvatarHandler)
	}

	return r
}
