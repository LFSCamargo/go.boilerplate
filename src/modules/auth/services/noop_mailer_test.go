package services

import (
	"context"
	"testing"

	"go.boilerplate/src/mail"
)

func TestNoopMailer_ReturnsNil(t *testing.T) {
	var mailer NoopMailer
	if err := mailer.SendVerifyEmail(context.Background(), mail.VerifyEmailParams{}); err != nil {
		t.Fatal(err)
	}
	if err := mailer.SendPasswordReset(context.Background(), mail.PasswordResetParams{}); err != nil {
		t.Fatal(err)
	}
}
