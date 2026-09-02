package config

import (
	"testing"
)

func TestNewConfig_RequiredEnv(t *testing.T) {
	t.Setenv("POSTGRES_CONNECTION", "postgres://user:pass@localhost:5432/test?sslmode=disable")
	t.Setenv("BEARER_SECRET", "test-secret-key-with-at-least-32-chars")
	t.Setenv("CORS_ORIGINS", "http://localhost:3000")

	cfg := NewConfig()

	if cfg.PostgresConnection == "" {
		t.Fatal("expected PostgresConnection to be set")
	}
	if cfg.BearerSecret == "" {
		t.Fatal("expected BearerSecret to be set")
	}
	if cfg.Port == "" {
		t.Fatal("expected Port default to be set")
	}
}
