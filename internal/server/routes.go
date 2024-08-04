package server

import (
	"Ultra-learn/internal/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/google/callback"),

	)

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/protected", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}), s.ProtectedHandler)
	v1 := r.Group("/api/v1")
	ar := v1.Group("/auth")
	ur := v1.Group("/user", middleware.AuthMiddleware, middleware.AccessControlMiddleware([]string{"admin", "user"}))
	{

		ar.POST("/register", s.RegisterUserHandler)
		ar.POST("/login", s.SignInUserHandler)
		ar.GET("/auth/:provider", s.SocialauthHandler)
		ar.GET("/:provider/callback", s.authCallback)
		ur.GET("/details", s.GetUserDetailsHandler)
		ur.PUT("/details", s.UpdateUserDetailsHandler)
		ur.POST("/avatar", s.UpdateAvatarHandler)
	}

	return r
}
