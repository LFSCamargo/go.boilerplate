package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.boilerplate/src/config"

	"github.com/gofiber/fiber/v3"
)

func TestAuthRateLimit_BlocksAfterMax(t *testing.T) {
	cfg := &config.Config{RateLimitMax: 2, RateLimitWindowMinutes: 1}
	app := fiber.New()
	app.Use(AuthRateLimit(cfg))
	app.Post("/api/v1/auth/login", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(`{}`)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("request %d: status %d", i+1, resp.StatusCode)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(`{}`)))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", resp.StatusCode)
	}
}

func TestAuthRateLimit_SkipsUnlistedRoutes(t *testing.T) {
	cfg := &config.Config{RateLimitMax: 1, RateLimitWindowMinutes: 1}
	app := fiber.New()
	app.Use(AuthRateLimit(cfg))
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("request %d: status %d", i+1, resp.StatusCode)
		}
	}
}
