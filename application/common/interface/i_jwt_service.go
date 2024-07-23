package jwt

import (
	"github.com/beka-birhanu/finance-go/domain/models"
	"github.com/dgrijalva/jwt-go"
)

type IJwtService interface {
	GenerateToken(user *models.User) (string, error)
	DecodeToken(token string) (jwt.MapClaims, error)
}
