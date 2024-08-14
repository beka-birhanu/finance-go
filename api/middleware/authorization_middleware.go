// Package middleware provides HTTP middleware for authorization and context handling.
package middleware

import (
	"context"
	"errors"
	"net/http"

	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
)

// contextKey is a type for context keys used in this package.
type contextKey string

// ContextUserClaims is the key for storing user claims in the context.
const ContextUserClaims contextKey = "userClaims"

// Authorization is a middleware that validates the JWT token from the request cookie.
// If the token is valid, the user claims are attached to the request context; otherwise,
// it returns an HTTP 401 Unauthorized error.
func Authorization(jwtService ijwt.IService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					http.Error(w, "Authorization token required", http.StatusUnauthorized)
				} else {
					http.Error(w, "Server error", http.StatusInternalServerError)
				}
				return
			}

			tokenString := cookie.Value
			claims, err := jwtService.Decode(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
