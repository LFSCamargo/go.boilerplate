package services

import (
	"context"
	"strings"
	"testing"
)

func newTestService(t *testing.T) (*AuthService, *fakeUsers, *fakeOTPs, *captureMailer) {
	t.Helper()
	users := newFakeUsers()
	otps := &fakeOTPs{}
	mailer := &captureMailer{}
	svc := NewAuthService(users, otps, newFakeRevoked(), mailer, Config{
		BearerSecret: "dev-secret-key-at-least-32-characters-long",
		AppPublicURL: "http://localhost:8080",
	})
	return svc, users, otps, mailer
}

func TestRegister_SendsVerification(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	user, err := svc.Register(context.Background(), RegisterInput{
		Email:       "Ada@Example.com",
		Password:    "password12",
		DisplayName: "Ada",
	})
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "ada@example.com" || user.EmailVerified {
		t.Fatalf("unexpected user: %+v", user)
	}
	if len(mailer.verify) != 1 || mailer.verify[0].Code == "" {
		t.Fatalf("expected verify email, got %+v", mailer.verify)
	}
	if !strings.Contains(mailer.verify[0].VerifyURL, "token=") {
		t.Fatalf("expected magic link, got %s", mailer.verify[0].VerifyURL)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	in := RegisterInput{Email: "ada@example.com", Password: "password12"}
	if _, err := svc.Register(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Register(context.Background(), in); err != ErrEmailTaken {
		t.Fatalf("got %v", err)
	}
}

func TestLogin_UnverifiedSendsEmail(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	mailer.verify = nil

	user, token, err := svc.Login(context.Background(), "ada@example.com", "password12")
	if err != nil {
		t.Fatalf("got %v", err)
	}
	if user == nil || token == nil || token.AccessToken == "" {
		t.Fatal("expected session for unverified user")
	}
	if user.EmailVerified {
		t.Fatal("expected unverified user")
	}
	if len(mailer.verify) != 1 {
		t.Fatal("expected verification email on login")
	}
}

func TestVerifyEmailCode_ThenLogin(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	code := mailer.verify[0].Code
	if err := svc.VerifyEmailCode(context.Background(), "ada@example.com", code); err != nil {
		t.Fatal(err)
	}

	user, token, err := svc.Login(context.Background(), "ada@example.com", "password12")
	if err != nil {
		t.Fatal(err)
	}
	if !user.EmailVerified || token == nil || token.AccessToken == "" {
		t.Fatalf("expected verified session: %+v %+v", user, token)
	}
}

func TestVerifyEmailToken(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	raw := mailer.verify[0].VerifyURL
	token := raw[strings.LastIndex(raw, "token=")+6:]
	if err := svc.VerifyEmailToken(context.Background(), token); err != nil {
		t.Fatal(err)
	}
	if _, session, err := svc.Login(context.Background(), "ada@example.com", "password12"); err != nil || session == nil {
		t.Fatalf("login after token verify: %v", err)
	}
}

func TestRecoverAndResetPassword(t *testing.T) {
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
	if len(mailer.reset) != 1 {
		t.Fatal("expected reset email")
	}
	if err := svc.ResetPassword(context.Background(), "ada@example.com", mailer.reset[0].Code, "newpass123"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Login(context.Background(), "ada@example.com", "password12"); err != ErrInvalidCredentials {
		t.Fatalf("old password should fail, got %v", err)
	}
	if _, token, err := svc.Login(context.Background(), "ada@example.com", "newpass123"); err != nil || token == nil {
		t.Fatalf("new password login: %v", err)
	}
}

func TestResetPasswordToken(t *testing.T) {
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
	raw := mailer.reset[0].ResetURL
	token := raw[strings.LastIndex(raw, "token=")+6:]
	if err := svc.ResetPasswordToken(context.Background(), token, "tokenpass1"); err != nil {
		t.Fatal(err)
	}
	if _, session, err := svc.Login(context.Background(), "ada@example.com", "tokenpass1"); err != nil || session == nil {
		t.Fatalf("token reset login: %v", err)
	}
}

func TestAuthenticate_LoadsUserAndRejectsRevoked(t *testing.T) {
	svc, _, _, mailer := newTestService(t)
	_, err := svc.Register(context.Background(), RegisterInput{Email: "ada@example.com", Password: "password12"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.VerifyEmailCode(context.Background(), "ada@example.com", mailer.verify[0].Code); err != nil {
		t.Fatal(err)
	}
	_, session, err := svc.Login(context.Background(), "ada@example.com", "password12")
	if err != nil {
		t.Fatal(err)
	}

	user, claims, err := svc.Authenticate(context.Background(), session.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "ada@example.com" || claims == nil || claims.Subject == "" {
		t.Fatalf("unexpected auth result: %+v %+v", user, claims)
	}

	if err := svc.Logout(context.Background(), session.AccessToken); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Authenticate(context.Background(), session.AccessToken); err != ErrInvalidToken {
		t.Fatalf("expected revoked token, got %v", err)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	svc, _, _, _ := newTestService(t)
	if _, err := svc.Register(context.Background(), RegisterInput{Email: "a@b.c", Password: "short"}); err != ErrWeakPassword {
		t.Fatalf("got %v", err)
	}
}
