package log

import (
	"strings"
	"testing"
)

func TestRedactValue_SensitiveKeys(t *testing.T) {
	tests := []struct {
		key  string
		in   string
		want string
	}{
		{"password", "password12", redacted},
		{"access_token", "abc123", redacted},
		{"authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.abc.def", "Bearer " + redacted},
		{"email", "ada@example.com", "a***@example.com"},
		{"postgres_connection", "postgres://goboilerplate:goboilerplate@localhost:5432/goboilerplate", "postgres://[REDACTED]:[REDACTED]@localhost:5432/goboilerplate"},
	}

	for _, tc := range tests {
		got := RedactValue(tc.key, tc.in)
		if got != tc.want {
			t.Fatalf("RedactValue(%q, %q) = %q, want %q", tc.key, tc.in, got, tc.want)
		}
	}
}

func TestRedactValue_NestedMap(t *testing.T) {
	in := map[string]any{
		"email":    "ada@example.com",
		"password": "secret123",
		"meta": map[string]any{
			"token": "abc",
		},
	}
	out := RedactValue("body", in).(map[string]any)
	if out["password"] != redacted {
		t.Fatalf("password not redacted: %#v", out["password"])
	}
	if out["email"] != "a***@example.com" {
		t.Fatalf("email not masked: %#v", out["email"])
	}
	meta := out["meta"].(map[string]any)
	if meta["token"] != redacted {
		t.Fatalf("nested token not redacted: %#v", meta["token"])
	}
}

func TestRedactQuery(t *testing.T) {
	got := RedactQuery("/api/v1/auth/verify-email?token=super-secret&foo=bar")
	if got == "" {
		t.Fatal("expected redacted query")
	}
	if strings.Contains(got, "super-secret") {
		t.Fatalf("token leaked in query: %s", got)
	}
	if !strings.Contains(got, "foo=bar") {
		t.Fatalf("non-sensitive query param removed: %s", got)
	}
}
