package openapi

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"go.boilerplate/src/middleware"
)

// Attach mounts Huma OpenAPI routes on the Fiber app (`/docs`, `/openapi.json`).
func Attach(app *fiber.App) huma.API {
	cfg := huma.DefaultConfig("Go Boilerplate API", "1.0.0")
	cfg.Info.Description = "Go Boilerplate REST API. Use **Authorize** in the playground with a Bearer JWT from `POST /api/v1/auth/login`."
	cfg.CreateHooks = nil
	cfg.Transformers = nil
	cfg.DocsPath = "/docs"
	cfg.OpenAPIPath = "/openapi"
	cfg.DocsRenderer = huma.DocsRendererScalar
	cfg.DocsRendererConfig = map[string]any{
		"theme":              "kepler",
		"darkMode":           true,
		"forceDarkModeState": "dark",
		"hideDarkModeToggle": true,
	}
	cfg.Servers = []*huma.Server{{
		URL:         "http://localhost:8080",
		Description: "Local development",
	}}
	if cfg.Components == nil {
		cfg.Components = &huma.Components{}
	}
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humafiber.New(app, cfg)
	api.UseMiddleware(middleware.AttachFiber)
	return api
}
