package tests

import (
	"Ultra-learn/internal/database"
	"Ultra-learn/internal/repository"
	"Ultra-learn/internal/server"
	"Ultra-learn/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func setupServer() *server.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbInstance := database.New()
	userRepo := repository.NewUserRepository(dbInstance.Db) // Pass the dbInstance to the UserRepository
	authService := services.NewAuthService(userRepo)        // Pass the UserRepository to the AuthService
	userService := services.NewUserService(userRepo)        // Assuming you have a function to create a new AuthService
	return &server.Server{
		Port:        port,
		Db:          dbInstance,
		AuthService: authService,
		UserService: userService,
	}
}

func TestHelloWorldHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.GET("/", s.HelloWorldHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"Hello World\",\"status_code\":200,\"data\":{\"message\":\"Hello World\"}}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestProtectedHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.GET("/protected", s.ProtectedHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"Hello World\",\"status_code\":200,\"data\":{\"message\":\"Hello World\"}}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRegisterUserHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.POST("/api/v1/auth/register", s.RegisterUserHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("POST", "/api/v1/auth/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	// Check the response body
	expected := "{\"message\":\"User created successfully\",\"status_code\":201}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSignInUserHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.POST("/api/v1/auth/login", s.SignInUserHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("POST", "/api/v1/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"User logged in successfully\",\"status_code\":200,\"data\":null}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetUserDetailsHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.GET("/api/v1/user/details", s.GetUserDetailsHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/v1/user/details", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"User details\",\"status_code\":200,\"data\":nul},"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
