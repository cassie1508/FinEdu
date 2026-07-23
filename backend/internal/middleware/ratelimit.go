package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit throttles requests per client IP using a token bucket, refilling
// at ratePerMinute/60 tokens per second with a burst of ratePerMinute. It's
// meant for unauthenticated, abuse-prone endpoints like signup.
func RateLimit(ratePerMinute int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)

	limit := rate.Limit(float64(ratePerMinute) / 60)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, ok := limiters[ip]
		if !ok {
			limiter = rate.NewLimiter(limit, ratePerMinute)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, please try again later"})
			return
		}

		c.Next()
	}
}
