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

	r.GET("/protected", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}), s.ProtectedHandler)

	v1 := r.Group("/api/v1")
	ar := v1.Group("/auth")
	ur := v1.Group("/user", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}))
	{

		ar.POST("/register", s.RegisterUserHandler)
		ar.POST("/login", s.SignInUserHandler)
		ur.GET("/details", s.GetUserDetailsHandler)
		ur.PUT("/details", s.UpdateUserDetailsHandler)
		ur.POST("/avatar", s.UpdateAvatarHandler)
	}

	return r
}
