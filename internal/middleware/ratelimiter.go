package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters  = make(map[string]*rate.Limiter)
	mu        sync.Mutex
	rateLimit = rate.Every(1 * time.Second) // 1 request per second
	burst     = 3                           // allow bursts of up to 3
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burst)
		limiters[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// Remove port if exists
	if host, _, err := net.SplitHostPort(ip); err == nil {
		return host
	}
	return ip
}
