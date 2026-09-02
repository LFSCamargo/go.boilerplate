package services

import (
	"context"
	"testing"
)

func TestPasswordReset_AllowsNewLogin(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.VerifyEmailCode(context.Background(), "ada@example.com", mailer.verify[0].Code); err != nil {
		t.Fatal(err)
	}
	if err := svc.RecoverPassword(context.Background(), "ada@example.com"); err != nil {
		t.Fatal(err)
	}
	if len(mailer.reset) != 1 || mailer.reset[0].Code == "" {
		t.Fatalf("expected reset email, got %+v", mailer.reset)
	}
	if err := svc.ResetPassword(context.Background(), "ada@example.com", mailer.reset[0].Code, "newpassword1"); err != nil {
		t.Fatal(err)
	}
	if _, token, err := svc.Login(context.Background(), "ada@example.com", "newpassword1"); err != nil || token == nil {
		t.Fatalf("login after reset: token=%v err=%v", token, err)
	}
}
