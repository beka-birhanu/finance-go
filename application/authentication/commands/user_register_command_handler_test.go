package commands

import (
	"testing"

	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/models"
	"github.com/dgrijalva/jwt-go"
)

type MockUserRepository struct {
	CreateUserFunc func(user *models.User) error
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserRepository) GetUserById(id string) (*models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) ListUser() ([]*models.User, error) {
	return nil, nil
}

type MockJwtService struct {
	GenerateTokenFunc func(user *models.User) (string, error)
}

func (m *MockJwtService) GenerateToken(user *models.User) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockJwtService) DecodeToken(token string) (jwt.MapClaims, error) {
	return nil, nil
}

type MockHashService struct {
}

func (m *MockHashService) Hash(word string) (string, error) {
	return "hashed" + word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return hashedWord == "hashed"+plainWord, nil
}

func TestUserRegisterCommandHandler_Handle(t *testing.T) {
	mockUserRepository := &MockUserRepository{
		CreateUserFunc: func(user *models.User) error {
			if user.Username() != "uniqueUsername" {
				return domain_errors.ErrUsernameConflict
			}
			return nil
		},
	}

	mockJwtService := &MockJwtService{
		GenerateTokenFunc: func(user *models.User) (string, error) {
			return "validToken", nil
		},
	}

	mockHashService := &MockHashService{}

	handler := NewRegisterCommandHandler(mockUserRepository, mockJwtService, mockHashService)

	validCommand, err := NewUserRegisterCommand("uniqueUsername", "#%strongPassword#%")
	if err != nil {
		t.Errorf("unexpected error '%v' on creating validCommand", err)
	}

	duplicateCommand, err := NewUserRegisterCommand("duplicateUsername", "#%strongPassword#%")
	if err != nil {
		t.Errorf("unexpected error '%v' on creating duplicateCommand", err)
	}

	weakPasswordCommand, err := NewUserRegisterCommand("uniqueUsername", "weakPassword")
	if err != nil {
		t.Errorf("unexpected error '%v' on creating weakPasswordCommand", err)
	}

	tests := []struct {
		name          string
		command       *UserRegisterCommand
		expectedError error
	}{
		{
			name:          "valid register",
			command:       validCommand,
			expectedError: nil,
		},
		{
			name:          "duplicate register",
			command:       duplicateCommand,
			expectedError: domain_errors.ErrUsernameConflict,
		},
		{
			name:          "weak password register",
			command:       weakPasswordCommand,
			expectedError: domain_errors.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(tt.command)
			if err != nil && tt.expectedError == nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
