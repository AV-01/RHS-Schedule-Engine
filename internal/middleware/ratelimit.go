package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(3 * time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 5*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute/60), 60)
	clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func ResetClients() {
	mu.Lock()
	defer mu.Unlock()
	clients = make(map[string]*client)
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getClient(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded, max 60 requests per minute.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
