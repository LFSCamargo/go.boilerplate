package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIssueAndParseAccessToken(t *testing.T) {
	secret := "dev-secret-key-at-least-32-characters-long"
	id := uuid.New()

	token, jti, _, err := IssueAccessToken(secret, id, "ada@example.com", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if jti == "" {
		t.Fatal("expected jti")
	}

	claims, err := ParseAccessToken(secret, token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != id.String() || claims.Email != "ada@example.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseAccessToken_RejectsBadSecret(t *testing.T) {
	id := uuid.New()
	token, _, _, err := IssueAccessToken("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", id, "a@b.c", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseAccessToken("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", token); err == nil {
		t.Fatal("expected signature error")
	}
}
