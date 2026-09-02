package healthRouter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.boilerplate/src/openapi"

	"github.com/gofiber/fiber/v3"
)

func TestHealthRouter(t *testing.T) {
	app := fiber.New()
	Register(openapi.Attach(app))

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to hit /health route: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"status":"ok"}`
	if strings.TrimSpace(string(body)) != expected {
		t.Errorf("Expected body %s, got %s", expected, string(body))
	}
}
