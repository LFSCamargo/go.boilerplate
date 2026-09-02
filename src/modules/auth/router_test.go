package authRouter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authHandlers "go.boilerplate/src/modules/auth/handlers"
	"go.boilerplate/src/openapi"

	"github.com/gofiber/fiber/v3"
)

var httpTest = fiber.TestConfig{Timeout: 10 * time.Second}

func TestAuthRouter_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	svc := newRouterService(t)
	Register(openapi.Attach(app), authHandlers.New(svc), svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("login route not registered")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/verify-email?token=bad", http.NoBody)
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("verify-email GET not registered")
	}
}

func TestAuthRouter_RegisterAndUnverifiedLogin(t *testing.T) {
	app := fiber.New()
	svc := newRouterService(t)
	Register(openapi.Attach(app), authHandlers.New(svc), svc)

	body, _ := json.Marshal(map[string]string{
		"email":    "router@example.com",
		"password": "password12",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("register status %d", resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("login unverified status %d", resp.StatusCode)
	}

	var loginResp struct {
		User struct {
			EmailVerified bool `json:"email_verified"`
		} `json:"user"`
		Token struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	if loginResp.User.EmailVerified || loginResp.Token.AccessToken == "" {
		t.Fatalf("expected unverified session: %+v", loginResp)
	}
}

func TestAuthRouter_VerifyEndpointsRequireAuth(t *testing.T) {
	app := fiber.New()
	svc := newRouterService(t)
	Register(openapi.Attach(app), authHandlers.New(svc), svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("resend verify status %d", resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email-code", bytes.NewBufferString(`{"code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("verify code status %d", resp.StatusCode)
	}
}

func TestAuthRouter_MeRequiresTokenAndReturnsUser(t *testing.T) {
	app := fiber.New()
	svc, users := newRouterBundle(t)
	Register(openapi.Attach(app), authHandlers.New(svc), svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", http.NoBody)
	resp, err := app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", resp.StatusCode)
	}

	body, _ := json.Marshal(map[string]string{
		"email":    "me@example.com",
		"password": "password12",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	users.markVerified("me@example.com")

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("login status %d", resp.StatusCode)
	}

	var loginResp struct {
		Token struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
	resp, err = app.Test(req, httpTest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("me status %d", resp.StatusCode)
	}

	var meResp struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&meResp); err != nil {
		t.Fatal(err)
	}
	if meResp.User.Email != "me@example.com" {
		t.Fatalf("unexpected me payload: %+v", meResp)
	}
}
