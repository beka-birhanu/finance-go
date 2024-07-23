package jwt

import (
	"testing"
	"time"

	"github.com/beka-birhanu/finance-go/domain/model"
)

type MockHashService struct {
	MatchFunc func(hashedWord, plainWord string) (bool, error)
}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return m.MatchFunc(hashedWord, plainWord)
}

var user, _ = model.NewUser("validUser", "#%@@strong@@password#%", &MockHashService{})

func TestJwtService(t *testing.T) {
	secretKey := "secret"
	issuer := "test_issuer"
	expTime := time.Minute * 15

	jwtService := NewJwtService(secretKey, issuer, expTime)

	t.Run("GenerateToken", func(t *testing.T) {
		token, err := jwtService.GenerateToken(user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token to be not empty")
		}
	})

	t.Run("DecodeToken", func(t *testing.T) {
		token, err := jwtService.GenerateToken(user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token to be not empty")
		}

		claims, err := jwtService.DecodeToken(token)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if claims["user_id"] != user.ID().String() {
			t.Errorf("expected user_id to be %v, got %v", user.ID().String(), claims["user_id"])
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
		_, err := jwtService.DecodeToken("invalid.token.string")
		if err == nil {
			t.Error("expected an error, got none")
		}
	})
}
