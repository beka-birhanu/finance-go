// Package jwt provides JWT generation and validation services.
package jwt

import (
	"errors"
	"time"

	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/dgrijalva/jwt-go"
)

// Service implements the ijwt.IService interface for handling JWT operations.
type Service struct {
	secretKey   string
	issuer      string
	expTime     time.Duration
	timeService itimeservice.IService
}

var _ ijwt.IService = &Service{}

// Config holds the configuration for creating a new JWT Service.
type Config struct {
	SecretKey   string
	Issuer      string
	ExpTime     time.Duration
	TimeService itimeservice.IService
}

// New creates a new JWT Service with the given configuration.
func New(config Config) *Service {
	return &Service{
		secretKey:   config.SecretKey,
		issuer:      config.Issuer,
		expTime:     config.ExpTime,
		timeService: config.TimeService,
	}
}

// Generate creates a new JWT token for the given user.
func (s *Service) Generate(user *usermodel.User) (string, error) {
	expirationTime := s.timeService.NowUTC().Add(s.expTime).Unix()
	claims := jwt.MapClaims{
		"user_id": user.ID().String(),
		"exp":     expirationTime,
		"iss":     s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// Decode parses and validates a JWT token, returning its claims if valid.
func (s *Service) Decode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, s.getSigningKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// getSigningKey is a helper function to validate the token's signing method and provide the secret key.
func (s *Service) getSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(s.secretKey), nil
}
