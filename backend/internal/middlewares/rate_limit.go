package middlewares

import (
	"sync"
	"time"

	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
}

// RateLimitMiddleware limits the number of requests per IP
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mutex.Lock()
		defer limiter.mutex.Unlock()

		now := time.Now()

		// Initialize if not exists
		if _, exists := limiter.requests[ip]; !exists {
			limiter.requests[ip] = []time.Time{}
		}

		// Remove old requests outside the window
		requests := limiter.requests[ip]
		validRequests := []time.Time{}
		for _, reqTime := range requests {
			if now.Sub(reqTime) < window {
				validRequests = append(validRequests, reqTime)
			}
		}

		// Check if limit exceeded
		if len(validRequests) >= maxRequests {
			utils.ErrorResponse(c, 429, "Rate limit exceeded", "Too many requests")
			c.Abort()
			return
		}

		// Add current request
		validRequests = append(validRequests, now)
		limiter.requests[ip] = validRequests

		c.Next()
	}
}

// APIKeyRateLimitMiddleware limits requests per API key
func APIKeyRateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.Next()
			return
		}

		limiter.mutex.Lock()
		defer limiter.mutex.Unlock()

		now := time.Now()
		key := "api_key_" + apiKey

		// Initialize if not exists
		if _, exists := limiter.requests[key]; !exists {
			limiter.requests[key] = []time.Time{}
		}

		// Remove old requests outside the window
		requests := limiter.requests[key]
		validRequests := []time.Time{}
		for _, reqTime := range requests {
			if now.Sub(reqTime) < window {
				validRequests = append(validRequests, reqTime)
			}
		}

		// Check if limit exceeded
		if len(validRequests) >= maxRequests {
			utils.ErrorResponse(c, 429, "Rate limit exceeded", "Too many requests for this API key")
			c.Abort()
			return
		}

		// Add current request
		validRequests = append(validRequests, now)
		limiter.requests[key] = validRequests

		c.Next()
	}
}
