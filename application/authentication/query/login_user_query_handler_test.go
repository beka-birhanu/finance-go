package query

import (
	"testing"
	"time"

	jwtInterface "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Mock implementations

type MockUserRepository struct {
	GetUserByUsernameFunc func(username string) (*model.User, error)
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	return nil
}

func (m *MockUserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	return m.GetUserByUsernameFunc(username)
}

func (m *MockUserRepository) ListUser() ([]*model.User, error) {
	return nil, nil
}

var _ repository.IUserRepository = &MockUserRepository{}

type MockJwtService struct {
	GenerateTokenFunc func(user *model.User) (string, error)
}

func (m *MockJwtService) GenerateToken(user *model.User) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockJwtService) DecodeToken(token string) (jwt.MapClaims, error) {
	return nil, nil
}

var _ jwtInterface.IJwtService = &MockJwtService{}

type MockHashService struct {
	MatchFunc func(hashedWord, plainWord string) (bool, error)
}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return m.MatchFunc(hashedWord, plainWord)
}

var _ hash.IHashService = &MockHashService{}

var validUser, _ = model.NewUser("validUser", "#%@@strong@@password#%", &MockHashService{}, time.Now().UTC())

func TestUserLoginQueryHandler_Handle(t *testing.T) {
	mockUserRepository := &MockUserRepository{
		GetUserByUsernameFunc: func(username string) (*model.User, error) {
			if username == "validUser" {
				return validUser, nil
			}
			return nil, appError.ErrUserNotFound
		},
	}

	mockJwtService := &MockJwtService{
		GenerateTokenFunc: func(user *model.User) (string, error) {
			return "validToken", nil
		},
	}

	mockHashService := &MockHashService{
		MatchFunc: func(hashedWord, plainWord string) (bool, error) {
			if hashedWord == validUser.PasswordHash() && plainWord == "password" {
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
			expectedError: appError.ErrInvalidUsernameOrPassword,
		},
		{
			name: "invalid password",
			query: &UserLoginQuery{
				Username: "validUser",
				Password: "wrongPassword",
			},
			expectedError: appError.ErrInvalidUsernameOrPassword,
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