package services

import (
	"context"
	"testing"

	"go.boilerplate/src/db/models"
)

func TestIssueVerification_SendsCodeAndLink(t *testing.T) {
	svc, users, _, mailer := newTestService(t)
	user, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	stored, err := users.FindByID(context.Background(), user.ID)
	if err != nil || stored == nil {
		t.Fatal(err)
	}
	mailer.verify = nil
	if err := svc.issueVerification(context.Background(), stored, users.policy); err != nil {
		t.Fatal(err)
	}
	if len(mailer.verify) != 1 || mailer.verify[0].Code == "" {
		t.Fatalf("got %+v", mailer.verify)
	}
	if mailer.verify[0].VerifyURL == "" {
		t.Fatal("expected verify URL")
	}
	if models.OTPPurposeEmailVerify == "" {
		t.Fatal("missing purpose")
	}
}
