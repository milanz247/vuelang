package middleware

import (
	"github.com/gin-gonic/gin"

	"vuelang/internal/framework/ratelimit"
	"vuelang/internal/framework/response"
)

// RateLimit returns a middleware that enforces the given Limiter per client IP.
func RateLimit(limiter *ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			response.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
