package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/application/common/interface/jwt"
)

type contextKey string

const (
	ContextUserClaims contextKey = "userClaims"
)

func AuthorizationMiddleware(jwtService jwt.IJwtService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("accessToken")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					http.Error(w, "Authorization token required", http.StatusUnauthorized)
				default:
					log.Println(err)
					http.Error(w, "server error", http.StatusInternalServerError)
				}
				return
			}

			tokenString := cookie.Value
			claims, err := jwtService.DecodeToken(tokenString)
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
