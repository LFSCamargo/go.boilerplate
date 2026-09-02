package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestSetup_ServesOpenAPIDocs(t *testing.T) {
	app := fiber.New()
	Setup(app, Deps{})

	req := httptest.NewRequest(http.MethodGet, "/docs", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("GET /docs status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	page := strings.ToLower(string(body))
	if !strings.Contains(page, "scalar") && !strings.Contains(page, "api-reference") {
		t.Fatalf("docs page missing playground markup: %s", body[:min(200, len(body))])
	}
}

func TestSetup_ServesOpenAPISpec(t *testing.T) {
	app := fiber.New()
	Setup(app, Deps{})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("GET /openapi.json status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	raw := string(body)
	if !strings.Contains(raw, `"openapi"`) {
		t.Fatalf("expected OpenAPI document, got %s", raw[:min(200, len(raw))])
	}
	if !strings.Contains(raw, `/health`) {
		t.Fatal("spec should document /health")
	}
}
