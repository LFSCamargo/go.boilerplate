package middleware

import (
	"net/http"
	"time"

	"go.boilerplate/src/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

var authRateLimitedRoutes = map[string]struct{}{
	http.MethodPost + " /api/v1/auth/login":             {},
	http.MethodPost + " /api/v1/auth/register":          {},
	http.MethodPost + " /api/v1/auth/recover-password":  {},
	http.MethodPost + " /api/v1/auth/reset-password":    {},
	http.MethodPost + " /api/v1/auth/verify-email":      {},
	http.MethodPost + " /api/v1/auth/verify-email-code": {},
	http.MethodGet + " /api/v1/auth/verify-email":       {},
}

// AuthRateLimit applies per-IP fixed-window limits to auth abuse targets.
func AuthRateLimit(cfg *config.Config) fiber.Handler {
	max := 10
	window := time.Minute
	if cfg != nil {
		if cfg.RateLimitMax > 0 {
			max = cfg.RateLimitMax
		}
		if cfg.RateLimitWindowMinutes > 0 {
			window = time.Duration(cfg.RateLimitWindowMinutes) * time.Minute
		}
	}

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		Next: func(c fiber.Ctx) bool {
			key := c.Method() + " " + c.Path()
			_, limited := authRateLimitedRoutes[key]
			return !limited
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate_limit_exceeded",
			})
		},
	})
}
