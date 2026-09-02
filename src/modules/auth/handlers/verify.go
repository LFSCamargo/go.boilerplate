package authHandlers

import (
	"context"
	"net/http"

	"go.boilerplate/src/middleware"
	"go.boilerplate/src/modules/auth/types/requests"
	"go.boilerplate/src/modules/auth/types/responses"
	"go.boilerplate/src/openapi"
)

type VerifyTokenInput struct {
	Token string `query:"token" doc:"Magic-link token from the verification email" minLength:"1"`
}

type VerifyCodeInput struct {
	Body requests.VerifyEmailCode
}

type VerifiedOutput struct {
	Body responses.Verified
}

func (h *Handlers) ResendVerification(ctx context.Context, _ *struct{}) (*OKOutput, error) {
	user, ok := middleware.ContextUser(ctx)
	if !ok {
		return nil, openapi.Status(http.StatusUnauthorized, "unauthorized")
	}
	if err := h.auth.ResendVerification(ctx, user.Email); err != nil {
		return nil, writeAuthError(err)
	}
	out := responses.OK{OK: true}
	if err := openapi.ValidateResponse(responses.OKSchema, out); err != nil {
		return nil, err
	}
	return &OKOutput{Body: out}, nil
}

func (h *Handlers) VerifyEmailToken(ctx context.Context, in *VerifyTokenInput) (*VerifiedOutput, error) {
	query := requests.VerifyEmailToken{Token: in.Token}
	if err := openapi.Validate(requests.VerifyEmailTokenSchema, query); err != nil {
		return nil, err
	}
	if err := h.auth.VerifyEmailToken(ctx, query.Token); err != nil {
		return nil, writeAuthError(err)
	}
	out := responses.Verified{OK: true, Verified: true}
	if err := openapi.ValidateResponse(responses.VerifiedSchema, out); err != nil {
		return nil, err
	}
	return &VerifiedOutput{Body: out}, nil
}

func (h *Handlers) VerifyEmailCode(ctx context.Context, in *VerifyCodeInput) (*VerifiedOutput, error) {
	if err := openapi.Validate(requests.VerifyEmailCodeSchema, in.Body); err != nil {
		return nil, err
	}
	user, ok := middleware.ContextUser(ctx)
	if !ok {
		return nil, openapi.Status(http.StatusUnauthorized, "unauthorized")
	}
	if err := h.auth.VerifyEmailCode(ctx, user.Email, in.Body.Code); err != nil {
		return nil, writeAuthError(err)
	}
	out := responses.Verified{OK: true, Verified: true}
	if err := openapi.ValidateResponse(responses.VerifiedSchema, out); err != nil {
		return nil, err
	}
	return &VerifiedOutput{Body: out}, nil
}
