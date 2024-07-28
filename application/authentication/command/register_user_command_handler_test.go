package registercmd

import (
	"errors"
	"testing"
	"time"

	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Mock implementations of interfaces for testing
type MockUserRepository struct {
	CreateUserFunc func(user *usermodel.User) error
}

func (m *MockUserRepository) Save(user *usermodel.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserRepository) Update(user *usermodel.User) error {
	return nil
}

func (m *MockUserRepository) ById(id uuid.UUID) (*usermodel.User, error) {
	return nil, nil
}

func (m *MockUserRepository) ByUsername(username string) (*usermodel.User, error) {
	return nil, nil
}

type MockJwtService struct {
	GenerateTokenFunc func(user *usermodel.User) (string, error)
}

func (m *MockJwtService) Generate(user *usermodel.User) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockJwtService) Decode(token string) (jwt.MapClaims, error) {
	return nil, nil
}

type MockHashService struct{}

func (m *MockHashService) Hash(word string) (string, error) {
	return "hashed" + word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return hashedWord == "hashed"+plainWord, nil
}

type MockTimeService struct{}

func (m *MockTimeService) NowUTC() time.Time {
	return time.Now().UTC()
}

// TestUserRegisterCommandHandler_Handle tests the Handle method of UserRegisterCommandHandler
func TestUserRegisterCommandHandler_Handle(t *testing.T) {
	mockUserRepository := &MockUserRepository{
		CreateUserFunc: func(user *usermodel.User) error {
			if user.Username() != "uniqueUsername" {
				return erruser.UsernameConflict
			}
			return nil
		},
	}

	mockJwtService := &MockJwtService{
		GenerateTokenFunc: func(user *usermodel.User) (string, error) {
			return "validToken", nil
		},
	}

	mockHashService := &MockHashService{}
	mockTimeService := &MockTimeService{}

	handler := NewHandler(Config{
		UserRepository: mockUserRepository,
		JwtService:     mockJwtService,
		HashService:    mockHashService,
		TimeService:    mockTimeService,
	})

	validCommand := &Command{Username: "uniqueUsername", Password: "#%strongPassword#%"}
	duplicateCommand := &Command{Username: "duplicateUsername", Password: "#%strongPassword#%"}
	weakPasswordCommand := &Command{Username: "uniqueUsername", Password: "weakPassword"}

	// Define the test cases
	tests := []struct {
		name          string
		command       *Command
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
			expectedError: erruser.UsernameConflict,
		},
		{
			name:          "weak password register",
			command:       weakPasswordCommand,
			expectedError: erruser.WeakPassword,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(tt.command)
			if err != nil && tt.expectedError == nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
			}

			if err != nil && tt.expectedError != nil {
				var unwrappedErr = errors.Unwrap(err).(*errdmn.Error)
				if unwrappedErr != nil && unwrappedErr != tt.expectedError {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			}
		})
	}
}
