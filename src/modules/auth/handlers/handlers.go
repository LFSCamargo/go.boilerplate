package authHandlers

import (
	"go.boilerplate/src/modules/auth/services"
	"go.boilerplate/src/openapi"
)

type Handlers struct {
	auth *services.AuthService
}

func New(auth *services.AuthService) *Handlers {
	return &Handlers{auth: auth}
}

func writeAuthError(err error) error {
	return openapi.FromAuth(err)
}
