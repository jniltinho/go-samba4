package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

type rateLimiter struct {
	sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var loginRatelimiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	limit:    5,
	window:   5 * time.Minute,
}

// RateLimit limits login attempts by IP to maxReqs per window
func RateLimit() echo.MiddlewareFunc {
	// Background cleanup of old visitors
	go func() {
		for {
			time.Sleep(time.Minute)
			loginRatelimiter.Lock()
			for ip, v := range loginRatelimiter.visitors {
				if time.Since(v.lastSeen) > loginRatelimiter.window {
					delete(loginRatelimiter.visitors, ip)
				}
			}
			loginRatelimiter.Unlock()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()

			loginRatelimiter.Lock()
			v, exists := loginRatelimiter.visitors[ip]
			if !exists {
				loginRatelimiter.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
				loginRatelimiter.Unlock()
				return next(c)
			}

			if time.Since(v.lastSeen) > loginRatelimiter.window {
				v.count = 1
				v.lastSeen = time.Now()
				loginRatelimiter.Unlock()
				return next(c)
			}

			if v.count >= loginRatelimiter.limit {
				loginRatelimiter.Unlock()
				return echo.NewHTTPError(http.StatusTooManyRequests, "Too many login attempts. Please try again later.")
			}

			v.count++
			v.lastSeen = time.Now()
			loginRatelimiter.Unlock()

			return next(c)
		}
	}
}
