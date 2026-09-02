package services

import (
	"context"

	"go.boilerplate/src/mail"
)

type NoopMailer struct{}

func (NoopMailer) SendVerifyEmail(context.Context, mail.VerifyEmailParams) error {
	return nil
}

func (NoopMailer) SendPasswordReset(context.Context, mail.PasswordResetParams) error {
	return nil
}
