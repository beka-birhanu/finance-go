package jwt

import (
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/dgrijalva/jwt-go"
)

type IJwtService interface {
	GenerateToken(user *model.User) (string, error)
	DecodeToken(token string) (jwt.MapClaims, error)
}
