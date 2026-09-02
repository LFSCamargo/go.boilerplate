package routes

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go.boilerplate/src/config"

	"github.com/gofiber/fiber/v3"
)

func TestSetup_RegistersHealthRoute(t *testing.T) {
	app := fiber.New()
	Setup(app, Deps{})

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSetup_ServesEmailStaticAssets(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := wd
	for {
		if _, err := os.Stat(filepath.Join(root, "emails", "static", "go_boilerplate", "logo_mark.svg")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Skip("emails/static/go_boilerplate not found")
		}
		root = parent
	}

	app := fiber.New()
	Setup(app, Deps{Config: &config.Config{
		EmailsDir:  filepath.Join(root, "emails"),
		UploadsDir: t.TempDir(),
	}})

	req := httptest.NewRequest(http.MethodGet, "/static/go_boilerplate/logo_mark.svg", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected static asset 200, got %d", resp.StatusCode)
	}
}
