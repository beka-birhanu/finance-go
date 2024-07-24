package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beka-birhanu/finance-go/api/user/dto"
	"github.com/beka-birhanu/finance-go/application/authentication/command"
	authCommand "github.com/beka-birhanu/finance-go/application/authentication/command"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/authentication/query"
	authQuery "github.com/beka-birhanu/finance-go/application/authentication/query"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Mock implementations for the IUserRepository interface
type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(user *model.User) error {
	return nil
}

func (m *mockUserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) ListUser() ([]*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) SomeRepoMethod() error {
	return nil
}

var _ repository.IUserRepository = &mockUserRepository{}

// Mock implementations for the IUserRegisterCommandHandler interface
type mockUserRegisterCommandHandler struct {
	handleFunc func(cmd *command.UserRegisterCommand) (*common.AuthResult, error)
}

func (m *mockUserRegisterCommandHandler) Handle(cmd *command.UserRegisterCommand) (*common.AuthResult, error) {
	return m.handleFunc(cmd)
}

var _ handlerInterface.ICommandHandler[*authCommand.UserRegisterCommand, *common.AuthResult] = &mockUserRegisterCommandHandler{}

// Mock implementations for the IUserLoginQueryHandler interface
type mockUserLoginQueryHandler struct {
	handleFunc func(query *query.UserLoginQuery) (*common.AuthResult, error)
}

func (m *mockUserLoginQueryHandler) Handle(q *query.UserLoginQuery) (*common.AuthResult, error) {
	return m.handleFunc(q)
}

var _ handlerInterface.ICommandHandler[*authQuery.UserLoginQuery, *common.AuthResult] = &mockUserLoginQueryHandler{}

func TestHandler_UserRegistrationAndLogin(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockRegisterCommandHandler := &mockUserRegisterCommandHandler{
		handleFunc: func(cmd *command.UserRegisterCommand) (*common.AuthResult, error) {
			switch cmd.Username {
			case "existinguser":
				return &common.AuthResult{}, appError.ErrUsernameConflict
			case "toolongusername":
				return &common.AuthResult{}, appError.ErrUsernameTooLong
			case "short":
				return &common.AuthResult{}, appError.ErrUsernameTooShort
			case "invalidformat!":
				return &common.AuthResult{}, appError.ErrUsernameInvalidFormat
			}
			if cmd.Password == "weakpassword" {
				return &common.AuthResult{}, appError.ErrWeakPassword
			}
			return &common.AuthResult{Token: "testtoken"}, nil
		},
	}
	mockLoginQueryHandler := &mockUserLoginQueryHandler{
		handleFunc: func(query *authQuery.UserLoginQuery) (*common.AuthResult, error) {
			if query.Username == "nonexistentuser" {
				return &common.AuthResult{}, appError.ErrInvalidUsernameOrPassword
			}
			if query.Password != "correctpassword" {
				return &common.AuthResult{}, appError.ErrInvalidUsernameOrPassword
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
			expectedError:  appError.ErrUsernameConflict.Message,
		},
		{
			name:           "Weak Password",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "newuser", Password: "weakpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  appError.ErrWeakPassword.Message,
		},
		{
			name:           "Username Too Long",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "toolongusername", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  appError.ErrUsernameTooLong.Message,
		},
		{
			name:           "Username Too Short",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "short", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  appError.ErrUsernameTooShort.Message,
		},
		{
			name:           "Username Invalid Format",
			url:            "/users/register",
			requestBody:    dto.RegisterRequest{Username: "invalidformat!", Password: "StrongPassword!123"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  appError.ErrUsernameInvalidFormat.Message,
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
			expectedError:  "invalid credentials",
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
