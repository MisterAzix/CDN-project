package app

import (
	"net/http"
	"golang.org/x/time/rate"
)

func NewRateLimiter(r rate.Limit, b int) *rate.Limiter {
	return rate.NewLimiter(r, b)
}

func RateLimitMiddleware(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}