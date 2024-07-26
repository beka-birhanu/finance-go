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
//
// Parameters:
//   - config: The configuration containing the secret key, issuer, expiration time, and time service.
//
// Returns:
//   - *Service: A new instance of the JWT Service.
func New(config Config) *Service {
	return &Service{
		secretKey:   config.SecretKey,
		issuer:      config.Issuer,
		expTime:     config.ExpTime,
		timeService: config.TimeService,
	}
}

// Generate creates a new JWT token for the given user.
//
// The token contains the user ID in the claims.
//
// Parameters:
//   - user: The user for whom the token is being generated.
//
// Returns:
//   - string: The signed JWT token.
//   - error: An error if the token could not be generated.
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
//
// If the token is expired or signed with a different method or secret key, an error is returned.
//
// Parameters:
//   - tokenString: The JWT token as a string.
//
// Returns:
//   - jwt.MapClaims: The claims extracted from the token.
//   - error: An error if the token is invalid or cannot be parsed.
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

