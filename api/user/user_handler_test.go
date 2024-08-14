package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beka-birhanu/finance-go/api/user/dto"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	queryHandlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Mock implementations for the IUserRepository interface
type MockUserRepository struct {
	CreateUserFunc func(user *usermodel.User) error
}

func (m *MockUserRepository) Save(user *usermodel.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserRepository) ById(id uuid.UUID) (*usermodel.User, error) {
	return nil, nil
}

func (m *MockUserRepository) ByUsername(username string) (*usermodel.User, error) {
	return nil, nil
}

var _ irepository.IUserRepository = &MockUserRepository{}

// Mock implementations for the IUserRegisterCommandHandler interface
type mockUserRegisterCommandHandler struct {
	handleFunc func(cmd *registercmd.Command) (*auth.Result, error)
}

func (m *mockUserRegisterCommandHandler) Handle(cmd *registercmd.Command) (*auth.Result, error) {
	return m.handleFunc(cmd)
}

var _ handlerInterface.IHandler[*registercmd.Command, *auth.Result] = &mockUserRegisterCommandHandler{}

// Mock implementations for the IUserLoginQueryHandler interface
type mockUserLoginQueryHandler struct {
	handleFunc func(query *loginqry.Query) (*auth.Result, error)
}

func (m *mockUserLoginQueryHandler) Handle(q *loginqry.Query) (*auth.Result, error) {
	return m.handleFunc(q)
}

var _ queryHandlerInterface.IHandler[*loginqry.Query, *auth.Result] = &mockUserLoginQueryHandler{}

func TestHandler_UserRegistrationAndLogin(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockRegisterCommandHandler := &mockUserRegisterCommandHandler{
		handleFunc: func(cmd *registercmd.Command) (*auth.Result, error) {
			switch cmd.Username {
			case "existinguser":
				return nil, fmt.Errorf("failed to create new user: %w", erruser.UsernameConflict)

			case "toolongusername":
				return nil, fmt.Errorf("failed to create new user: %w", erruser.UsernameTooLong)
			case "short":
				return nil, fmt.Errorf("failed to create new user: %w", erruser.UsernameTooShort)
			case "invalidformat!":
				return nil, fmt.Errorf("failed to create new user: %w", erruser.UsernameInvalidFormat)
			}
			if cmd.Password == "weakpassword" {
				return nil, fmt.Errorf("failed to create new user: %w", erruser.WeakPassword)
			}
			return auth.NewResult(uuid.New(), cmd.Username, "testtoken"), nil
		},
	}
	mockLoginQueryHandler := &mockUserLoginQueryHandler{
		handleFunc: func(query *loginqry.Query) (*auth.Result, error) {
			if query.Username == "nonexistentuser" {
				return nil, appError.InvalidCredential("user does not exist")
			}
			if query.Password != "correctpassword" {
				return nil, appError.InvalidCredential("incorrect passoword")
			}
			return auth.NewResult(uuid.New(), query.Username, "testtoken"), nil
		},
	}

	config := Config{
		UserRepository:  mockRepo,
		RegisterHandler: mockRegisterCommandHandler,
		LoginHandler:    mockLoginQueryHandler,
	}
	h := NewHandler(config)

	router := mux.NewRouter()
	h.RegisterPublic(router)

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
			expectedError:  erruser.UsernameConflict.Error(),
		},
		{
			name:           "Weak Password",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "newuser", Password: "weakpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  erruser.WeakPassword.Error(),
		},
		{
			name:           "Username Too Long",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "toolongusername", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  erruser.UsernameTooLong.Error(),
		},
		{
			name:           "Username Too Short",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "short", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  erruser.UsernameTooShort.Error(),
		},
		{
			name:           "Username Invalid Format",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "invalidformat!", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  erruser.UsernameInvalidFormat.Error(),
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
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
		},
		{
			name:           "Invalid Credentials",
			url:            "/users/login",
			requestBody:    dto.LoginUserRequest{Username: "existinguser", Password: "wrongpassword"},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
		},
		{
			name:           "Invalid Login Request Body",
			url:            "/users/login",
			requestBody:    struct{}{},
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
