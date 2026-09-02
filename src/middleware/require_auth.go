package middleware

import (
	"context"

	"go.boilerplate/src/modules/auth/services"

	"github.com/gofiber/fiber/v3"
)

const (
	localUser  = "auth_user"
	localToken = "auth_token"
	localJTI   = "auth_jti"
)

// Authenticator is implemented by the auth service (and test doubles).
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*services.AuthUser, *services.Claims, error)
}

// RequireAuth validates the Bearer token, rejects revoked JWTs, and injects
// the authenticated user into the Fiber context for any module's handlers.
func RequireAuth(auth Authenticator) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := bearerToken(c)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing_token"})
		}

		user, claims, err := auth.Authenticate(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		c.Locals(localUser, user)
		c.Locals(localToken, token)
		c.Locals(localJTI, claims.ID)
		return c.Next()
	}
}

func User(c fiber.Ctx) (*services.AuthUser, bool) {
	user, ok := c.Locals(localUser).(*services.AuthUser)
	return user, ok && user != nil
}

func bearerToken(c fiber.Ctx) string {
	header := c.Get("Authorization")
	const prefix = "Bearer "
	if len(header) > len(prefix) && header[:len(prefix)] == prefix {
		return header[len(prefix):]
	}
	return ""
}
