package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestRedactingHandler_RedactsAttrs(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	h := &RedactingHandler{inner: inner}

	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "auth", 0)
	rec.AddAttrs(
		slog.String("email", "ada@example.com"),
		slog.String("password", "password12"),
		slog.String("authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.abc.def"),
	)
	if err := h.Handle(context.Background(), rec); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "password12") {
		t.Fatalf("password leaked: %s", out)
	}
	if strings.Contains(out, "ada@example.com") {
		t.Fatalf("raw email leaked: %s", out)
	}
	if !strings.Contains(out, "a***@example.com") {
		t.Fatalf("masked email missing: %s", out)
	}
	if strings.Contains(out, "eyJhbGciOi") {
		t.Fatalf("jwt leaked: %s", out)
	}
}

func TestParseLevel(t *testing.T) {
	if parseLevel("debug") != slog.LevelDebug {
		t.Fatal("expected debug level")
	}
	if parseLevel("unknown") != slog.LevelInfo {
		t.Fatal("expected default info level")
	}
}
