package tests

import (
	"Ultra-learn/internal/database"
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/logger"
	"Ultra-learn/internal/repository"
	"Ultra-learn/internal/server"
	"Ultra-learn/internal/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func LoadEnv() {
	err := godotenv.Load(".env") // adjust the path to your .env file if needed
	if err != nil {
		log.Fatalf("Error loading .env file" + err.Error())
	}
}

func SetupServer() *server.Server {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)

	}
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	log.Println("config", os.Getenv("DB_DATABASE"))
	var (
		databaseName = os.Getenv("DB_DATABASE")
		dbport       = os.Getenv("DB_PORT")
		dbhost       = os.Getenv("DB_HOST")
	)
	log.Println("PORT: ", port)
	dbInstance := database.New(databaseName, dbhost, dbport)
	logger.Init()
	userRepo := repository.NewUserRepository(dbInstance) // Pass the dbInstance to the UserRepository
	authService := services.NewAuthService(userRepo)     // Pass the UserRepository to the AuthService
	userService := services.NewUserService(userRepo)     // Assuming you have a function to create a new AuthService
	emailService := services.NewEmailService()
	return &server.Server{
		Port:         port,
		Db:           dbInstance,
		AuthService:  authService,
		UserService:  userService,
		EmailService: emailService,
	}
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
func AssertBool(t *testing.T, got, expected bool) {
	if got != expected {
		t.Errorf("handler returned wrong boolean: got %v expected %v", got, expected)
	}
}

func AssertValidationError(t *testing.T, response map[string]interface{}, field string, expectedMessage string) {
	errors, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'error' field in response")
	}

	errorMessage, exists := errors[field]
	if !exists {
		t.Fatalf("expected validation error message for field '%s'", field)
	}

	if errorMessage != expectedMessage {
		t.Errorf("unexpected error messagserver := setupServer()e for field '%s': got %v, want %v", field, errorMessage, expectedMessage)
	}
}

// helper to signup a user
func SignupUser(t *testing.T, r *gin.Engine, testServer *server.Server, userSignUpData dto.CreateUserRequest, admin bool) {

	var (
		signupPath = "/api/v1/auth/register"
		signupURI  = url.URL{Path: signupPath}
	)

	r.POST(signupPath, testServer.RegisterUserHandler)

	// if admin {
	// 	signupPath = "/api/v1/auth/admin/signup"
	// 	signupURI = url.URL{Path: signupPath}
	// 	r.POST(signupPath, auth.CreateAdmin)
	// }

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(userSignUpData)
	req, err := http.NewRequest(http.MethodPost, signupURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
}

func GetLoginToken(t *testing.T, r *gin.Engine, testServer *server.Server, loginData dto.LoginRequest) string {
	var (
		loginPath = "/api/v1/auth/login"
		loginURI  = url.URL{Path: loginPath}
	)
	r.POST(loginPath, testServer.SignInUserHandler)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(loginData)
	req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		return ""
	}

	data := ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	token := dataM["access_token"].(string)

	return token
}
