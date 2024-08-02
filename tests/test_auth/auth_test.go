package tests

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/helper"
	"Ultra-learn/tests"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserSignUp(t *testing.T) {
	s := tests.SetupServer()
	requestURI := url.URL{Path: "/api/v1/auth/signup"}
	randomID := helper.GenerateUserId()

	testCases := []struct {
		name     string
		payload  dto.CreateUserRequest
		expected int
		Message  string
	}{
		{
			name: "valid user signup",
			payload: dto.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     fmt.Sprintf("primejay166+%s@gmail.com", randomID),
				Password:  "password",
			},
			expected: 201,
			Message:  "User created successfully",
		},
		{
			name: "invalid email",
			payload: dto.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     fmt.Sprintf("primejay166+%s@gmail", randomID),
				Password:  "password",
			},
			expected: 400,
			Message:  "Invalid email address",
		},
		{
			name: "User already exists",
			payload: dto.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     fmt.Sprintf("primejay166@gmail.com"),
				Password:  "password",
			},
			expected: 400,
			Message:  "User already exists",
		},
	}
	for _, tt := range testCases {
		r := gin.Default()

		r.POST(requestURI.Path, s.RegisterUserHandler)
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			// Marshal the payload to json
			json.NewEncoder(&b).Encode(tt.payload)
			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			// Assert the status code
			tests.AssertStatusCode(t, rr.Code, tt.expected)
			data := tests.ParseResponse(rr)
			code := int(data["code"].(float64))
			tests.AssertStatusCode(t, code, tt.expected)
			if tt.Message != "" {
				message := data["message"]
				if message != nil {
					tests.AssertResponseMessage(t, message.(string), tt.Message)
				} else {
					tests.AssertResponseMessage(t, "", tt.Message)
				}

			}

		})
	}

}
