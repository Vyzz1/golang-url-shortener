package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	count     int
	lastReset time.Time
}

type rateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	rate    int
	window  time.Duration
}

func newRateLimiter(rate int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		clients: make(map[string]*client),
		rate:    rate,
		window:  window,
	}
	go rl.cleanupClients()
	return rl
}

func (rl *rateLimiter) cleanupClients() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastReset) > rl.window {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &client{
			count:     1,
			lastReset: time.Now(),
		}
		return true
	}

	if time.Since(c.lastReset) > rl.window {
		c.count = 1
		c.lastReset = time.Now()
		return true
	}

	if c.count < rl.rate {
		c.count++
		return true
	}

	return false
}

var limiter = newRateLimiter(60, time.Minute) // 60 requests per minute

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
