package log

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

func baseHandler(w io.Writer, level string, format string) slog.Handler {
	lvl := parseLevel(level)
	opts := &slog.HandlerOptions{Level: lvl}

	if strings.EqualFold(format, "json") {
		return slog.NewJSONHandler(w, opts)
	}

	return tint.NewHandler(w, &tint.Options{
		Level:      lvl,
		TimeFormat: time.Kitchen,
		NoColor:    os.Getenv("NO_COLOR") != "",
	})
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
