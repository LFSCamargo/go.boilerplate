package services

import (
	"context"
	"testing"
)

func TestResendVerification_AlreadyVerified(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.VerifyEmailCode(context.Background(), "ada@example.com", mailer.verify[0].Code); err != nil {
		t.Fatal(err)
	}
	if err := svc.ResendVerification(context.Background(), "ada@example.com"); err != ErrAlreadyVerified {
		t.Fatalf("got %v", err)
	}
}

func TestVerifyEmailToken_ConfirmsAccount(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	url := mailer.verify[0].VerifyURL
	token := url[len(url)-len(""):]
	idx := 0
	for i := range url {
		if url[i] == '=' {
			idx = i + 1
		}
	}
	token = url[idx:]
	if err := svc.VerifyEmailToken(context.Background(), token); err != nil {
		t.Fatal(err)
	}
	if _, session, err := svc.Login(context.Background(), "ada@example.com", "password12"); err != nil || session == nil {
		t.Fatalf("login after token verify: token=%v err=%v", session, err)
	}
}
