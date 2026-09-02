package openapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestAttach_ServesDocsAndSpec(t *testing.T) {
	app := fiber.New()
	Attach(app)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Go Boilerplate API") {
		t.Fatalf("spec title missing: %s", body[:min(300, len(body))])
	}
	if !strings.Contains(string(body), "localhost:8080") {
		t.Fatal("spec should list the local :8080 server")
	}

	docsReq := httptest.NewRequest(http.MethodGet, "/docs", http.NoBody)
	docsResp, err := app.Test(docsReq)
	if err != nil {
		t.Fatal(err)
	}
	defer docsResp.Body.Close()
	docsBody, _ := io.ReadAll(docsResp.Body)
	page := string(docsBody)
	if !strings.Contains(page, "scalar") && !strings.Contains(page, "api-reference") {
		t.Fatalf("docs missing Scalar: %s", page[:min(200, len(page))])
	}
	if !strings.Contains(page, "kepler") && !strings.Contains(page, "darkMode") {
		t.Fatal("docs should request Scalar dark mode")
	}
}
