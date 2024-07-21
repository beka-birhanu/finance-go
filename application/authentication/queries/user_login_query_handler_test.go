package queries

import (
	"errors"
	"testing"

	"github.com/beka-birhanu/finance-go/domain/entities"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Mock implementations

type MockUserRepository struct {
	GetUserByUsernameFunc func(username string) (*entities.User, error)
}

func (m *MockUserRepository) CreateUser(user *entities.User) error {
	return nil
}

func (m *MockUserRepository) GetUserById(id string) (*entities.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*entities.User, error) {
	return m.GetUserByUsernameFunc(username)
}

func (m *MockUserRepository) ListUser() ([]*entities.User, error) {
	return nil, nil
}

type MockJwtService struct {
	GenerateTokenFunc func(user *entities.User) (string, error)
}

func (m *MockJwtService) GenerateToken(user *entities.User) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockJwtService) DecodeToken(token string) (jwt.MapClaims, error) {
	return nil, nil
}

type MockHashService struct {
	MatchFunc func(hashedWord, plainWord string) (bool, error)
}

func (m *MockHashService) Hash(word string) (string, error) {
	return "", nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return m.MatchFunc(hashedWord, plainWord)
}

func TestUserLoginQueryHandler_Handle(t *testing.T) {
	mockUserRepository := &MockUserRepository{
		GetUserByUsernameFunc: func(username string) (*entities.User, error) {
			if username == "validUser" {
				return &entities.User{
					ID:       uuid.New(),
					Username: "validUser",
					Password: "hashedPassword",
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	mockJwtService := &MockJwtService{
		GenerateTokenFunc: func(user *entities.User) (string, error) {
			return "validToken", nil
		},
	}

	mockHashService := &MockHashService{
		MatchFunc: func(hashedWord, plainWord string) (bool, error) {
			if hashedWord == "hashedPassword" && plainWord == "password" {
				return true, nil
			}
			return false, nil
		},
	}

	handler := NewUserLoginQueryHandler(mockUserRepository, mockJwtService, mockHashService)

	tests := []struct {
		name          string
		query         *UserLoginQuery
		expectedError error
	}{
		{
			name: "valid login",
			query: &UserLoginQuery{
				Username: "validUser",
				Password: "password",
			},
			expectedError: nil,
		},
		{
			name: "invalid username",
			query: &UserLoginQuery{
				Username: "invalidUser",
				Password: "password",
			},
			expectedError: errors.New("invalid username or password"),
		},
		{
			name: "invalid password",
			query: &UserLoginQuery{
				Username: "validUser",
				Password: "wrongPassword",
			},
			expectedError: errors.New("invalid username or password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(tt.query)
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
