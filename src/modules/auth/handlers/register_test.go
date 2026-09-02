package authHandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/modules/auth/services"
	"go.boilerplate/src/openapi"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v3"
)

func testRegisterAPI(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	svc := services.NewAuthService(
		&memUsers{users: map[string]*models.User{}},
		&memOTPs{},
		memRevoked{},
		services.NoopMailer{},
		services.Config{BearerSecret: "dev-secret-key-at-least-32-characters-long", AppPublicURL: "http://localhost:8080"},
	)
	api := openapi.Attach(app)
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/register",
		DefaultStatus: http.StatusCreated,
	}, New(svc).Register)
	return app
}

func TestRegister_JSONCreated(t *testing.T) {
	app := testRegisterAPI(t)
	body, _ := json.Marshal(map[string]string{"email": "ada@example.com", "password": "password12", "display_name": "Ada"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, fiber.TestConfig{Timeout: 10 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestRegister_InvalidJSONBody(t *testing.T) {
	app := testRegisterAPI(t)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":"not-an-email"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, fiber.TestConfig{Timeout: 10 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected 400 or 422, got %d", resp.StatusCode)
	}
}
