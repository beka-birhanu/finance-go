package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beka-birhanu/finance-go/api/users/dto"
	"github.com/beka-birhanu/finance-go/application/authentication/commands"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/authentication/queries"
	commandAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	querieAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/models.go"
	"github.com/gorilla/mux"
)

// Mock implementations for the IUserRepository interface
type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(user *models.User) error {
	return nil
}

func (m *mockUserRepository) GetUserById(id string) (*models.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (m *mockUserRepository) ListUser() ([]*models.User, error) {
	return nil, nil
}

func (m *mockUserRepository) SomeRepoMethod() error {
	return nil
}

var _ persistance.IUserRepository = &mockUserRepository{}

// Mock implementations for the IUserRegisterCommandHandler interface
type mockUserRegisterCommandHandler struct {
	handleFunc func(cmd *commands.UserRegisterCommand) (*common.AuthResult, error)
}

func (m *mockUserRegisterCommandHandler) Handle(cmd *commands.UserRegisterCommand) (*common.AuthResult, error) {
	return m.handleFunc(cmd)
}

var _ commandAuth.IUserRegisterCommandHandler = &mockUserRegisterCommandHandler{}

// Mock implementations for the IUserLoginQueryHandler interface
type mockUserLoginQueryHandler struct {
	handleFunc func(query *queries.UserLoginQuery) (*common.AuthResult, error)
}

func (m *mockUserLoginQueryHandler) Handle(q *queries.UserLoginQuery) (*common.AuthResult, error) {
	return m.handleFunc(q)
}

var _ querieAuth.IUserLoginQueryHandler = &mockUserLoginQueryHandler{}

func TestHandler_UserRegistrationAndLogin(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockRegisterCommandHandler := &mockUserRegisterCommandHandler{
		handleFunc: func(cmd *commands.UserRegisterCommand) (*common.AuthResult, error) {
			if cmd.Username == "existinguser" {
				return &common.AuthResult{}, domain_errors.ErrUsernameConflict
			}
			if cmd.Password == "weakpassword" {
				return &common.AuthResult{}, domain_errors.ErrWeakPassword
			}
			return &common.AuthResult{Token: "testtoken"}, nil
		},
	}
	mockLoginQueryHandler := &mockUserLoginQueryHandler{
		handleFunc: func(query *queries.UserLoginQuery) (*common.AuthResult, error) {
			if query.Username == "nonexistentuser" {
				return &common.AuthResult{}, errors.New("user not found")
			}
			if query.Password != "correctpassword" {
				return &common.AuthResult{}, errors.New("invalid credentials")
			}
			return &common.AuthResult{Token: "testtoken"}, nil
		},
	}

	h := NewHandler(mockRepo, mockRegisterCommandHandler, mockLoginQueryHandler)

	router := mux.NewRouter()
	h.RegisterPublicRoutes(router)

	tests := []struct {
		name           string
		url            string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Successful Registration",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "newuser", Password: "StrongPassword!123"},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Username In Use",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "existinguser", Password: "StrongPassword!123"},
			expectedStatus: http.StatusConflict,
			expectedError:  domain_errors.ErrUsernameConflict.Error(),
		},
		{
			name:           "Weak Password",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "newuser", Password: "weakpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  domain_errors.ErrWeakPassword.Error(),
		},
		{
			name:           "Invalid Register Request Body",
			url:            "/users/register",
			requestBody:    struct{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid payload: Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:           "Successful Login",
			url:            "/users/login",
			requestBody:    dto.LoginUserRequest{Username: "existinguser", Password: "correctpassword"},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Nonexistent User",
			url:            "/users/login",
			requestBody:    dto.LoginUserRequest{Username: "nonexistentuser", Password: "correctpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "user not found",
		},
		{
			name:           "Invalid Credentials",
			url:            "/users/login",
			requestBody:    dto.LoginUserRequest{Username: "existinguser", Password: "wrongpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid credentials",
		},
		{
			name:           "Invalid Login Request Body",
			url:            "/users/login",
			requestBody:    struct{}{}, // empty object
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid payload: Key: 'LoginUserRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'LoginUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedError != "" {
				var response map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if response["error"] != tt.expectedError {
					t.Errorf("handler returned unexpected error: got %v want %v", response["error"], tt.expectedError)
				}
			}
		})
	}
}
