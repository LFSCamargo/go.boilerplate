package log

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestHTTP_LogsWithoutAuthorization(t *testing.T) {
	Init("debug", "json")

	var buf bytes.Buffer
	orig := defaultLogger
	t.Cleanup(func() { defaultLogger = orig })

	defaultLogger = slogWithWriter(&buf)

	app := fiber.New()
	app.Use(HTTP())
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	req.Header.Set("Authorization", "Bearer super-secret-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	out := buf.String()
	if strings.Contains(out, "super-secret-token") {
		t.Fatalf("authorization leaked into logs: %s", out)
	}
	if !strings.Contains(out, "/health") {
		t.Fatalf("expected path in logs: %s", out)
	}
}

func slogWithWriter(w io.Writer) *slog.Logger {
	return slog.New(newHandler(w, "debug", "json"))
}
