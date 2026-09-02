package authHandlers

import (
	"context"
	"net/http"

	"go.boilerplate/src/middleware"
	"go.boilerplate/src/modules/auth/types/responses"
	"go.boilerplate/src/openapi"
)

type MeOutput struct {
	Body responses.Me
}

func (h *Handlers) Me(ctx context.Context, _ *struct{}) (*MeOutput, error) {
	user, ok := middleware.ContextUser(ctx)
	if !ok {
		return nil, openapi.Status(http.StatusUnauthorized, "unauthorized")
	}
	out := responses.Me{User: responses.UserFrom(user)}
	if err := openapi.ValidateResponse(responses.MeSchema, out); err != nil {
		return nil, err
	}
	return &MeOutput{Body: out}, nil
}
