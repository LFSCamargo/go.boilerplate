package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.boilerplate/src/config"
	"go.boilerplate/src/db"
	"go.boilerplate/src/mail"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func integrationDSN() string {
	if v := os.Getenv("POSTGRES_CONNECTION"); v != "" {
		return v
	}
	return "postgres://goboilerplate:goboilerplate@localhost:5432/goboilerplate?sslmode=disable"
}

func setupIntegrationApp(t *testing.T) *fiber.App {
	t.Helper()
	dsn := integrationDSN()
	database, err := db.Connect(dsn)
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	if err := database.Ping(); err != nil {
		t.Skipf("postgres ping failed: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	if err := db.RunMigrations(dsn); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg := &config.Config{
		Port:               "5000",
		PostgresConnection: dsn,
		BearerSecret:       "dev-secret-key-at-least-32-characters-long",
		CORSOrigins:        "http://localhost:3000",
		SMTPHost:           "localhost",
		SMTPPort:           "1025",
		SMTPSecurity:       "none",
		MailFrom:           "noreply@go.boilerplate",
		EmailsDir:          findEmailsDir(t),
		EmailAssetsBaseURL: "http://localhost:8080",
		AppPublicURL:       "http://localhost:8080",
		UploadsDir:         t.TempDir(),
	}

	var mailer *mail.Service
	if svc, err := mail.NewFromConfig(cfg); err == nil {
		mailer = svc
	}

	app := fiber.New()
	Setup(app, Deps{Config: cfg, DB: database, Mail: mailer})
	return app
}

func findEmailsDir(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "emails", "package.json")); err == nil {
			return filepath.Join(wd, "emails")
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Skip("emails workspace not found")
		}
		wd = parent
	}
}

func TestIntegration_AuthFlow(t *testing.T) {
	app := setupIntegrationApp(t)
	email := fmt.Sprintf("test-auth-%s@example.com", uuid.NewString())

	registerBody, err := json.Marshal(map[string]string{
		"email":        email,
		"password":     "password12",
		"display_name": "Ada",
	})
	if err != nil {
		t.Fatal(err)
	}

	slow := fiber.TestConfig{Timeout: 30 * time.Second}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, slow)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("register: %d %s", resp.StatusCode, body)
	}

	loginBody, _ := json.Marshal(map[string]string{"email": email, "password": "password12"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, slow)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("unverified login: %d %s", resp.StatusCode, body)
	}

	var unverifiedLogin struct {
		Token struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&unverifiedLogin); err != nil {
		t.Fatal(err)
	}
	if unverifiedLogin.Token.AccessToken == "" {
		t.Fatal("expected token for unverified login")
	}

	code, token := latestMailhogSecrets(t, email)
	if token == "" && code == "" {
		t.Skip("mailhog not available; signup/login verified against API only")
	}
	if token == "" {
		t.Skip("mailhog HTML had no verify token (six-digit scrape is unreliable)")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/verify-email?token="+token, http.NoBody)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("magic link verify: %d %s", resp.StatusCode, body)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("verified login: %d %s", resp.StatusCode, body)
	}

	var loginResp struct {
		Token struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("logout: %d", resp.StatusCode)
	}
}

func latestMailhogSecrets(t *testing.T, email string) (code, token string) {
	t.Helper()
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", ""
	}
	var payload struct {
		Items []struct {
			Content struct {
				Headers map[string][]string `json:"Headers"`
				Body    string              `json:"Body"`
			} `json:"Content"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", ""
	}
	for i := len(payload.Items) - 1; i >= 0; i-- {
		item := payload.Items[i]
		to := strings.Join(item.Content.Headers["To"], ",")
		if !strings.Contains(to, email) {
			continue
		}
		body := item.Content.Body
		if idx := strings.Index(body, "token="); idx >= 0 {
			raw := body[idx+6:]
			token = strings.TrimRight(raw[:min(64, len(raw))], `"'<> `)
		}
		for _, n := range extractSixDigit(body) {
			code = n
		}
		return code, token
	}
	return "", ""
}

func extractSixDigit(s string) []string {
	var out []string
	for i := 0; i+6 <= len(s); i++ {
		ok := true
		for j := 0; j < 6; j++ {
			if s[i+j] < '0' || s[i+j] > '9' {
				ok = false
				break
			}
		}
		if ok && (i == 0 || s[i-1] < '0' || s[i-1] > '9') {
			out = append(out, s[i:i+6])
		}
	}
	return out
}
