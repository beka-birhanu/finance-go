package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
)

// contextKey defines a type for context keys used in the middleware.
type contextKey string

// Constants for context keys
const (
	// ContextUserClaims is the key used to store user claims in the context.
	ContextUserClaims contextKey = "userClaims"
)

// Authorization is a middleware function that validates the JWT token from the request cookie.
// It uses the provided JWT service to decode the token and attach the user claims to the request context.
// If the token is invalid or missing, it responds with an appropriate HTTP error.
func Authorization(jwtService ijwt.IService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the access token from the request cookie
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					http.Error(w, "Authorization token required", http.StatusUnauthorized)
				default:
					log.Println(err)
					http.Error(w, "Server error", http.StatusInternalServerError)
				}
				return
			}

			tokenString := cookie.Value
			// Decode the token using the JWT service
			claims, err := jwtService.Decode(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Attach claims to the context
			ctx := context.WithValue(r.Context(), ContextUserClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

