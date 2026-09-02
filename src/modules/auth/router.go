package authRouter

import (
	"net/http"

	"go.boilerplate/src/middleware"
	authHandlers "go.boilerplate/src/modules/auth/handlers"
	"go.boilerplate/src/modules/auth/services"

	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API, handlers *authHandlers.Handlers, auth *services.AuthService) {
	huma.Register(api, huma.Operation{
		OperationID:   "auth-login",
		Method:        http.MethodPost,
		Path:          "/api/v1/auth/login",
		Summary:       "Log in",
		Description:   "Issue a JWT. Unverified accounts receive a session with `email_verified: false` and a fresh verification email.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusOK,
	}, handlers.Login)

	huma.Register(api, huma.Operation{
		OperationID:   "auth-register",
		Method:        http.MethodPost,
		Path:          "/api/v1/auth/register",
		Summary:       "Register",
		Description:   "Create an account. Optional `picture` multipart field is accepted alongside JSON or form fields.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusCreated,
	}, handlers.Register)

	huma.Register(api, huma.Operation{
		OperationID: "auth-recover-password",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/recover-password",
		Summary:     "Request password reset",
		Tags:        []string{"Auth"},
	}, handlers.RecoverPassword)

	huma.Register(api, huma.Operation{
		OperationID: "auth-reset-password",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/reset-password",
		Summary:     "Reset password",
		Description: "Reset with OTP `code` + `email`, or a magic-link `token`.",
		Tags:        []string{"Auth"},
	}, handlers.ResetPassword)

	huma.Register(api, huma.Operation{
		OperationID: "auth-verify-email-token",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/verify-email",
		Summary:     "Verify email via magic link",
		Tags:        []string{"Auth"},
	}, handlers.VerifyEmailToken)

	protected := huma.NewGroup(api)
	protected.UseMiddleware(middleware.RequireHumaAuth(auth))

	huma.Register(protected, huma.Operation{
		OperationID: "auth-resend-verification",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/verify-email",
		Summary:     "Resend verification email",
		Description: "Requires Bearer JWT. Uses the authenticated user's email.",
		Tags:        []string{"Auth"},
		Security:    []map[string][]string{{"BearerAuth": {}}},
	}, handlers.ResendVerification)

	huma.Register(protected, huma.Operation{
		OperationID: "auth-verify-email-code",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/verify-email-code",
		Summary:     "Verify email with OTP",
		Description: "Requires Bearer JWT. Submits the 6-digit code for the authenticated user.",
		Tags:        []string{"Auth"},
		Security:    []map[string][]string{{"BearerAuth": {}}},
	}, handlers.VerifyEmailCode)

	huma.Register(protected, huma.Operation{
		OperationID: "auth-me",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/me",
		Summary:     "Current user",
		Tags:        []string{"Auth"},
		Security:    []map[string][]string{{"BearerAuth": {}}},
	}, handlers.Me)

	huma.Register(protected, huma.Operation{
		OperationID:   "auth-logout",
		Method:        http.MethodPost,
		Path:          "/api/v1/auth/logout",
		Summary:       "Log out",
		Description:   "Revoke the current access token `jti`.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusNoContent,
		Security:      []map[string][]string{{"BearerAuth": {}}},
	}, handlers.Logout)
}
