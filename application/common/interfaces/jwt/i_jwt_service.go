package jwt

import (
	"github.com/beka-birhanu/finance-go/domain/entities"
	"github.com/dgrijalva/jwt-go"
)

type IJwtService interface {
	GenerateToken(user *entities.User) (string, error)
	DecodeToken(token string) (jwt.MapClaims, error)
}

