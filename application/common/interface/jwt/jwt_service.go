/*
Package jwt provides an interface for handling JSON Web Tokens (JWT).

It includes the `IJwtService` interface for generating and decoding JWTs.
*/
package ijwt

import (
	"github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/dgrijalva/jwt-go"
)

// IJwtService defines methods for generating and decoding JSON Web Tokens (JWT).
//
// Methods:
// - GenerateToken(user *usermodel.User) (string, error): Generates a JWT for the
// given user and returns the token or an error.
// - DecodeToken(token string) (jwt.MapClaims, error): Decodes the provided JWT
// and returns the claims or an error.
type IJwtService interface {
	// GenerateToken creates a JWT for the specified user.
	GenerateToken(user *usermodel.User) (string, error)

	// DecodeToken parses the provided JWT and returns the claims or an error.
	DecodeToken(token string) (jwt.MapClaims, error)
}
