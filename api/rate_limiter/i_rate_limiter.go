package ratelimiter

import "golang.org/x/time/rate"

type IRateLimiter interface {
	AddLimiter(ip string) *rate.Limiter
	GetLimiter(ip string) *rate.Limiter
}
