package loginqry

import (
	"testing"
	"time"

	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Mock implementations

type MockUserRepository struct {
	ByUsernameFunc func(username string) (*usermodel.User, error)
	AddFunc        func(user *usermodel.User) error
}

func (m *MockUserRepository) Add(user *usermodel.User) error {
	return m.AddFunc(user)
}

func (m *MockUserRepository) Update(user *usermodel.User) error {
	return nil
}

func (m *MockUserRepository) ById(id uuid.UUID) (*usermodel.User, error) {
	return nil, nil
}

func (m *MockUserRepository) ByUsername(username string) (*usermodel.User, error) {
	return m.ByUsernameFunc(username)
}

var _ irepository.IUserRepository = &MockUserRepository{}

type MockJwtService struct {
	GenerateTokenFunc func(user *usermodel.User) (string, error)
}

func (m *MockJwtService) Generate(user *usermodel.User) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockJwtService) Decode(token string) (jwt.MapClaims, error) {
	return nil, nil
}

var _ ijwt.IJwtService = &MockJwtService{}

type MockHashService struct {
	MatchFunc func(hashedWord, plainWord string) (bool, error)
	HashFunc  func(word string) (string, error)
}

func (m *MockHashService) Hash(word string) (string, error) {
	return m.HashFunc(word)
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return m.MatchFunc(hashedWord, plainWord)
}

var _ hash.IHashService = &MockHashService{}

var validUser, _ = usermodel.New(usermodel.Config{
	Username:       "validUser",
	PlainPassword:  `#%@@strong@@password#%`,
	CreationTime:   time.Now().UTC(),
	PasswordHasher: &MockHashService{},
},
)

func TestHandler_Handle(t *testing.T) {
	mockUserRepository := &MockUserRepository{
		ByUsernameFunc: func(username string) (*usermodel.User, error) {
			if username == "validUser" {
				return validUser, nil
			}
			return nil, errdmn.NewNotFound("user not found")
		},
		AddFunc: func(user *usermodel.User) error {
			if user.Username() == "newUser" {
				return nil
			}
			return erruser.UsernameConflict
		},
	}

	mockJwtService := &MockJwtService{
		GenerateTokenFunc: func(user *usermodel.User) (string, error) {
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
		HashFunc: func(word string) (string, error) {
			return "#%@@strong@@password#%", nil
		},
	}

	handler := NewHandler(Config{
		UserRepository: mockUserRepository,
		JwtService:     mockJwtService,
		HashService:    mockHashService,
	})

	tests := []struct {
		name          string
		query         *Query
		expectedError error
	}{
		{
			name: "valid login",
			query: &Query{
				Username: "validUser",
				Password: "password",
			},
			expectedError: nil,
		},
		{
			name: "invalid username",
			query: &Query{
				Username: "invalidUser",
				Password: "password",
			},
			expectedError: appError.InvalidCredential("user not found"),
		},
		{
			name: "invalid password",
			query: &Query{
				Username: "validUser",
				Password: "wrongPassword",
			},
			expectedError: appError.InvalidCredential("incorrect password"),
		},
		{
			name: "new user registration",
			query: &Query{
				Username: "newUser",
				Password: "password",
			},
			expectedError: nil,
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

