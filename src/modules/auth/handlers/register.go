package authHandlers

import (
	"context"
	"net/http"

	"go.boilerplate/src/middleware"
	"go.boilerplate/src/modules/auth/services"
	"go.boilerplate/src/modules/auth/types/requests"
	"go.boilerplate/src/modules/auth/types/responses"
	"go.boilerplate/src/openapi"
)

type RegisterInput struct {
	Body requests.Register
}

type RegisterOutput struct {
	Status int `json:"-"`
	Body   responses.RegisterCreated
}

func (h *Handlers) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	if err := openapi.Validate(requests.RegisterSchema, in.Body); err != nil {
		return nil, err
	}
	payload := services.RegisterInput{
		Email:       in.Body.Email,
		Password:    in.Body.Password,
		DisplayName: in.Body.DisplayName,
	}
	if fc := middleware.ContextFiber(ctx); fc != nil {
		if file, err := fc.FormFile("picture"); err == nil {
			payload.Picture = file
		}
	}

	user, err := h.auth.Register(ctx, payload)
	if err != nil {
		return nil, writeAuthError(err)
	}

	out := responses.RegisterCreated{
		User:                 responses.UserFrom(user),
		RequiresVerification: true,
	}
	if err := openapi.ValidateResponse(responses.RegisterCreatedSchema, out); err != nil {
		return nil, err
	}
	return &RegisterOutput{Status: http.StatusCreated, Body: out}, nil
}
