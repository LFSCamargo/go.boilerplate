package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"go.boilerplate/src/modules/auth/services"
)

type denyAuth struct{}

func (denyAuth) Authenticate(context.Context, string) (*services.AuthUser, *services.Claims, error) {
	return nil, nil, services.ErrInvalidToken
}

func TestRequireHumaAuth_RejectsMissingToken(t *testing.T) {
	app := fiber.New()
	cfg := huma.DefaultConfig("t", "1")
	cfg.CreateHooks = nil
	cfg.Transformers = nil
	api := humafiber.New(app, cfg)
	grp := huma.NewGroup(api)
	grp.UseMiddleware(RequireHumaAuth(denyAuth{}))
	huma.Register(grp, huma.Operation{
		OperationID: "secret",
		Method:      http.MethodGet,
		Path:        "/secret",
	}, func(context.Context, *struct{}) (*struct {
		Body string `json:"ok"`
	}, error) {
		return &struct {
			Body string `json:"ok"`
		}{Body: "yes"}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/secret", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d %s", resp.StatusCode, body)
	}
}
