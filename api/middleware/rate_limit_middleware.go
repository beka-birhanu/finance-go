package middleware

import (
	"log"
	"net/http"

	ratelimiter "github.com/beka-birhanu/finance-go/api/rate_limiter"
)

func RateLimitMiddleware(limiter ratelimiter.IRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			ratelimiter := limiter.GetLimiter(ip)
			log.Println("serving", ip, ratelimiter.Allow(), ratelimiter.Limit())

			if !ratelimiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
