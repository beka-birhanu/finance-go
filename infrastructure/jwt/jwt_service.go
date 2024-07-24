package jwt

import (
	"errors"
	"time"

	jwtInterface "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/dgrijalva/jwt-go"
)

type JwtService struct {
	secretKey   string
	issuer      string
	expTime     time.Duration
	timeService timeservice.ITimeService
}

var _ jwtInterface.IJwtService = &JwtService{}

func NewJwtService(secretKey, issuer string, expTime time.Duration, timeService timeservice.ITimeService) *JwtService {
	return &JwtService{
		secretKey:   secretKey,
		issuer:      issuer,
		expTime:     expTime,
		timeService: timeService,
	}
}

func (s *JwtService) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID().String(),
		"exp":     s.timeService.NowUTC().Add(s.expTime).Unix(),
		"iss":     s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
