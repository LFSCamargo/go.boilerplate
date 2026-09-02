package log

import (
	"log/slog"
	"testing"
)

func TestInit_SetsDefaultLogger(t *testing.T) {
	Init("debug", "json")
	if Default() == nil {
		t.Fatal("expected default logger")
	}
}

func TestRedactArgs_Pairs(t *testing.T) {
	args := redactArgs("email", "ada@example.com", "password", "secret123")
	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d", len(args))
	}
	if args[1] != "a***@example.com" {
		t.Fatalf("email not masked: %#v", args[1])
	}
	if args[3] != redacted {
		t.Fatalf("password not redacted: %#v", args[3])
	}
}

func TestLogLevelMapping(t *testing.T) {
	if parseLevel("error") != slog.LevelError {
		t.Fatal("expected error level")
	}
}
