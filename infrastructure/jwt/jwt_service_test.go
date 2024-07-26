package jwt

import (
	"testing"
	"time"

	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
)

// MockHashService mocks the hash service for testing.
type MockHashService struct {
	MatchFunc func(hashedWord, plainWord string) (bool, error)
}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return m.MatchFunc(hashedWord, plainWord)
}

var _ hash.IService = &MockHashService{}

// MockTimeService mocks the time service for testing.
type MockTimeService struct{}

func (m *MockTimeService) NowUTC() time.Time {
	return time.Now().UTC()
}

var _ timeservice.IService = &MockTimeService{}

// Test user setup
var testUser, _ = usermodel.New(usermodel.Config{
	Username:       "validUser",
	PlainPassword:  "#%@@strong@@password#%",
	CreationTime:   time.Now().UTC(),
	PasswordHasher: &MockHashService{},
})

func TestJwtService(t *testing.T) {
	secretKey := "secret"
	issuer := "test_issuer"
	expTime := time.Minute * 15

	jwtService := New(Config{
		SecretKey:   secretKey,
		Issuer:      issuer,
		ExpTime:     expTime,
		TimeService: &MockTimeService{},
	})

	t.Run("GenerateToken", func(t *testing.T) {
		token, err := jwtService.Generate(testUser)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token to be not empty")
		}
	})

	t.Run("DecodeToken", func(t *testing.T) {
		token, err := jwtService.Generate(testUser)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token to be not empty")
		}

		claims, err := jwtService.Decode(token)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if claims["user_id"] != testUser.ID().String() {
			t.Errorf("expected user_id to be %v, got %v", testUser.ID().String(), claims["user_id"])
		}
		if claims["iss"] != issuer {
			t.Errorf("expected issuer to be %v, got %v", issuer, claims["iss"])
		}

		exp := int64(claims["exp"].(float64))
		if !time.Unix(exp, 0).After(time.Now()) {
			t.Error("expected exp to be in the future")
		}
	})

	t.Run("DecodeInvalidToken", func(t *testing.T) {
		_, err := jwtService.Decode("invalid.token.string")
		if err == nil {
			t.Error("expected an error, got none")
		}
	})
}

