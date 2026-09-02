package middleware

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/modules/auth/services"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type stubUsers struct {
	user *models.User
}

func (s *stubUsers) Create(context.Context, *models.User) error { return nil }
func (s *stubUsers) FindByEmail(context.Context, string) (*models.User, error) {
	return s.user, nil
}
func (s *stubUsers) FindByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	if s.user != nil && s.user.ID == id {
		clone := *s.user
		return &clone, nil
	}
	return nil, nil
}
func (s *stubUsers) Update(context.Context, *models.User) error { return nil }
func (s *stubUsers) RegistrationConfig(context.Context) (*models.UserRegistrationConfig, error) {
	return &models.UserRegistrationConfig{AllowRegistration: true, MinPasswordLength: 8}, nil
}

func (s *stubUsers) SaveAvatar(*multipart.FileHeader) (string, error) {
	return "", nil
}

type stubOTP struct{}

func (stubOTP) Create(context.Context, *models.OTP) error                         { return nil }
func (stubOTP) ConsumeActive(context.Context, uuid.UUID, models.OTPPurpose) error { return nil }
func (stubOTP) FindLatestActive(context.Context, uuid.UUID, models.OTPPurpose) (*models.OTP, error) {
	return nil, nil
}
func (stubOTP) FindActiveByTokenHash(context.Context, string) (*models.OTP, error) {
	return nil, nil
}
func (stubOTP) Save(context.Context, *models.OTP) error { return nil }

type stubRevoked struct {
	jtis map[string]bool
}

func (s *stubRevoked) Revoke(_ context.Context, _ uuid.UUID, jti string, _ time.Time) error {
	if s.jtis == nil {
		s.jtis = map[string]bool{}
	}
	s.jtis[jti] = true
	return nil
}

func (s *stubRevoked) IsRevoked(_ context.Context, jti string) (bool, error) {
	return s.jtis[jti], nil
}

func testAuthService(t *testing.T) (*services.AuthService, *models.User) {
	t.Helper()
	id := uuid.New()
	name := "Ada"
	user := &models.User{
		ID:            id,
		Email:         "ada@example.com",
		DisplayName:   &name,
		EmailVerified: true,
		PasswordHash:  "x",
	}
	svc := services.NewAuthService(
		&stubUsers{user: user},
		stubOTP{},
		&stubRevoked{jtis: map[string]bool{}},
		services.NoopMailer{},
		services.Config{BearerSecret: "dev-secret-key-at-least-32-characters-long"},
	)
	return svc, user
}

func TestRequireAuth_MissingToken(t *testing.T) {
	svc, _ := testAuthService(t)
	app := fiber.New()
	app.Get("/me", RequireAuth(svc), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestRequireAuth_InjectsUser(t *testing.T) {
	svc, user := testAuthService(t)
	token, _, _, err := services.IssueAccessToken("dev-secret-key-at-least-32-characters-long", user.ID, user.Email, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	app.Get("/me", RequireAuth(svc), func(c fiber.Ctx) error {
		got, ok := User(c)
		if !ok || got.Email != user.Email {
			t.Fatalf("user not injected: ok=%v user=%+v", ok, got)
		}
		return c.JSON(fiber.Map{"email": got.Email})
	})

	req := httptest.NewRequest(http.MethodGet, "/me", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status %d body %s", resp.StatusCode, body)
	}
}
