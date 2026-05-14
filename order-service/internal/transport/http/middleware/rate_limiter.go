package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimiter is the contract the middleware depends on.
type RateLimiter interface {
	Allow(ctx context.Context, clientID string) (bool, error)
}

// RateLimit returns a Gin middleware that enforces per-IP request limits.
// On Redis failure it fails open so legitimate traffic is never blocked by infra issues.
func RateLimit(limiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := limiter.Allow(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.Next()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded — try again later",
			})
			return
		}
		c.Next()
	}
}
