package ratelimiter

import (
	"sync"
	"time"

	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiters for different IPs with eviction support.
type IPRateLimiter struct {
	limiters    map[string]*LimiterEntry
	mu          sync.RWMutex
	rate        rate.Limit
	burst       int
	timeService itimeservice.IService
}

// Ensure IPRateLimiter implements IRateLimiter.
var _ IRateLimiter = &IPRateLimiter{}

// LimiterEntry contains the rate limiter and the last time it was accessed.
type LimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// NewIPRateLimiter initializes a new IPRateLimiter with the given rate and burst size.
func NewIPRateLimiter(r rate.Limit, b int, t itimeservice.IService) *IPRateLimiter {
	i := &IPRateLimiter{
		limiters:    make(map[string]*LimiterEntry),
		rate:        r,
		burst:       b,
		timeService: t,
	}
	go i.cleanupExpiredEntries()
	return i
}

// AddLimiter creates a new rate limiter for the provided IP and adds it to the map.
func (i *IPRateLimiter) AddLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.limiters[ip] = &LimiterEntry{
		limiter:    limiter,
		lastAccess: i.timeService.NowUTC(),
	}

	return limiter
}

// GetLimiter retrieves the rate limiter for the provided IP. If it doesn't exist, a new one is created.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	entry, exists := i.limiters[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddLimiter(ip)
	}

	i.mu.Lock()
	entry.lastAccess = i.timeService.NowUTC()
	i.mu.Unlock()

	return entry.limiter
}

// cleanupExpiredEntries periodically removes rate limiters that haven't been used for a certain amount of time.
func (i *IPRateLimiter) cleanupExpiredEntries() {
	expTime := 10 * time.Minute
	for {
		time.Sleep(expTime)
		i.mu.Lock()

		for ip, entry := range i.limiters {
			if time.Since(entry.lastAccess) > expTime {
				delete(i.limiters, ip)
			}
		}

		i.mu.Unlock()
	}
}
