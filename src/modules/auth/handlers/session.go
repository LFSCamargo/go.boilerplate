package authHandlers

import (
	"context"
	"net/http"

	"go.boilerplate/src/middleware"
	"go.boilerplate/src/modules/auth/types/requests"
	"go.boilerplate/src/modules/auth/types/responses"
	"go.boilerplate/src/openapi"
)

type LoginInput struct {
	Body requests.Login
}

type SessionOutput struct {
	Body responses.Session
}

type OKOutput struct {
	Body responses.OK
}

type RecoverInput struct {
	Body requests.Email
}

type ResetPasswordInput struct {
	Body requests.ResetPassword
}

func (h *Handlers) Login(ctx context.Context, in *LoginInput) (*SessionOutput, error) {
	if err := openapi.Validate(requests.LoginSchema, in.Body); err != nil {
		return nil, err
	}
	user, token, err := h.auth.Login(ctx, in.Body.Email, in.Body.Password)
	if err != nil {
		return nil, writeAuthError(err)
	}

	out := responses.Session{
		User:  responses.UserFrom(user),
		Token: responses.TokenFrom(token),
	}
	if err := openapi.ValidateResponse(responses.SessionSchema, out); err != nil {
		return nil, err
	}
	return &SessionOutput{Body: out}, nil
}

func (h *Handlers) Logout(ctx context.Context, _ *struct{}) (*struct{}, error) {
	if _, ok := middleware.ContextUser(ctx); !ok {
		return nil, openapi.Status(http.StatusUnauthorized, "unauthorized")
	}
	if err := h.auth.Logout(ctx, middleware.ContextToken(ctx)); err != nil {
		return nil, writeAuthError(err)
	}
	return nil, nil
}

func (h *Handlers) RecoverPassword(ctx context.Context, in *RecoverInput) (*OKOutput, error) {
	if err := openapi.Validate(requests.EmailSchema, in.Body); err != nil {
		return nil, err
	}
	if err := h.auth.RecoverPassword(ctx, in.Body.Email); err != nil {
		return nil, writeAuthError(err)
	}
	out := responses.OK{OK: true}
	if err := openapi.ValidateResponse(responses.OKSchema, out); err != nil {
		return nil, err
	}
	return &OKOutput{Body: out}, nil
}

func (h *Handlers) ResetPassword(ctx context.Context, in *ResetPasswordInput) (*OKOutput, error) {
	if err := openapi.Validate(requests.ResetPasswordSchema, in.Body); err != nil {
		return nil, err
	}
	if !in.Body.HasToken() && !in.Body.HasCode() {
		return nil, openapi.Status(http.StatusBadRequest, "invalid_body")
	}

	var err error
	if in.Body.HasToken() {
		err = h.auth.ResetPasswordToken(ctx, in.Body.Token, in.Body.NewPassword)
	} else {
		err = h.auth.ResetPassword(ctx, in.Body.Email, in.Body.Code, in.Body.NewPassword)
	}
	if err != nil {
		return nil, writeAuthError(err)
	}
	out := responses.OK{OK: true}
	if err := openapi.ValidateResponse(responses.OKSchema, out); err != nil {
		return nil, err
	}
	return &OKOutput{Body: out}, nil
}
