package log

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// HTTP returns Fiber middleware that logs requests with PII-safe fields.
func HTTP() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if status == 0 {
			status = fiber.StatusOK
		}

		path := RedactQuery(c.OriginalURL())
		attrs := []any{
			"method", c.Method(),
			"path", path,
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", c.IP(),
		}
		if rid := c.Get("X-Request-ID"); rid != "" {
			attrs = append(attrs, "request_id", rid)
		}
		if ua := c.Get("User-Agent"); ua != "" {
			attrs = append(attrs, "user_agent", ua)
		}

		switch {
		case status >= 500:
			Error("request", attrs...)
		case status >= 400:
			Warn("request", attrs...)
		default:
			Info("request", attrs...)
		}

		return err
	}
}
