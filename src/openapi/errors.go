package openapi

import (
	"errors"
	"net/http"

	"go.boilerplate/src/modules/auth/services"
)

// APIError is a compact JSON error: `{"error":"..."}`.
type APIError struct {
	status int
	Code   string `json:"error"`
}

func (e *APIError) GetStatus() int { return e.status }
func (e *APIError) Error() string  { return e.Code }

func Status(code int, codeName string) error {
	return &APIError{status: code, Code: codeName}
}

func FromAuth(err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		return Status(http.StatusUnauthorized, "invalid_credentials")
	case errors.Is(err, services.ErrEmailTaken):
		return Status(http.StatusConflict, "email_taken")
	case errors.Is(err, services.ErrEmailNotVerified):
		return Status(http.StatusForbidden, "email_not_verified")
	case errors.Is(err, services.ErrInvalidOTP), errors.Is(err, services.ErrInvalidToken):
		return Status(http.StatusBadRequest, "invalid_or_expired")
	case errors.Is(err, services.ErrRegistrationClosed):
		return Status(http.StatusForbidden, "registration_closed")
	case errors.Is(err, services.ErrWeakPassword):
		return Status(http.StatusBadRequest, "weak_password")
	case errors.Is(err, services.ErrInvalidPicture):
		return Status(http.StatusBadRequest, "invalid_picture")
	case errors.Is(err, services.ErrAlreadyVerified):
		return Status(http.StatusConflict, "already_verified")
	default:
		return Status(http.StatusInternalServerError, "internal_error")
	}
}
