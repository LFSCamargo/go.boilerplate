package middleware

import (
	"context"
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"go.boilerplate/src/modules/auth/services"
)

type ctxKey string

const (
	ctxFiber ctxKey = "fiber_ctx"
	ctxUser  ctxKey = localUser
	ctxToken ctxKey = localToken
)

// AttachFiber stores the Fiber context on the request context for FormFile, etc.
func AttachFiber(ctx huma.Context, next func(huma.Context)) {
	next(huma.WithValue(ctx, ctxFiber, humafiber.Unwrap(ctx)))
}

// RequireHumaAuth is the Huma equivalent of RequireAuth.
func RequireHumaAuth(auth Authenticator) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token := bearerFromHeader(ctx.Header("Authorization"))
		if token == "" {
			writeHumaJSON(ctx, fiber.StatusUnauthorized, fiber.Map{"error": "missing_token"})
			return
		}
		user, claims, err := auth.Authenticate(ctx.Context(), token)
		if err != nil {
			writeHumaJSON(ctx, fiber.StatusUnauthorized, fiber.Map{"error": "unauthorized"})
			return
		}
		ctx = huma.WithValue(ctx, ctxUser, user)
		ctx = huma.WithValue(ctx, ctxToken, token)
		ctx = huma.WithValue(ctx, ctxKey(localJTI), claims.ID)
		next(ctx)
	}
}

func ContextUser(ctx context.Context) (*services.AuthUser, bool) {
	user, ok := ctx.Value(ctxUser).(*services.AuthUser)
	return user, ok && user != nil
}

func ContextToken(ctx context.Context) string {
	token, _ := ctx.Value(ctxToken).(string)
	return token
}

func ContextFiber(ctx context.Context) fiber.Ctx {
	c, _ := ctx.Value(ctxFiber).(fiber.Ctx)
	return c
}

func bearerFromHeader(header string) string {
	const prefix = "Bearer "
	if len(header) > len(prefix) && header[:len(prefix)] == prefix {
		return header[len(prefix):]
	}
	return ""
}

func writeHumaJSON(ctx huma.Context, status int, payload any) {
	ctx.SetStatus(status)
	ctx.SetHeader("Content-Type", "application/json")
	_ = json.NewEncoder(ctx.BodyWriter()).Encode(payload)
}
