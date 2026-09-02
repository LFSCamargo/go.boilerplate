package log

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	defaultLogger *slog.Logger
	initOnce      sync.Once
)

// Init configures the global logger. format is "pretty" (dev) or "json" (prod).
func Init(level string, format string) {
	initOnce.Do(func() {
		defaultLogger = slog.New(newHandler(os.Stdout, level, format))
		slog.SetDefault(defaultLogger)
	})
}

// InitFromEnv reads LOG_LEVEL and LOG_FORMAT from the environment.
func InitFromEnv() {
	Init(envOr("LOG_LEVEL", "info"), envOr("LOG_FORMAT", "pretty"))
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// Default returns the configured application logger.
func Default() *slog.Logger {
	if defaultLogger == nil {
		InitFromEnv()
	}
	return defaultLogger
}

func Info(msg string, args ...any)  { Default().Info(msg, redactArgs(args...)...) }
func Warn(msg string, args ...any)  { Default().Warn(msg, redactArgs(args...)...) }
func Error(msg string, args ...any) { Default().Error(msg, redactArgs(args...)...) }

func Fatal(msg string, args ...any) {
	Default().Error(msg, redactArgs(args...)...)
	os.Exit(1)
}

func redactArgs(args ...any) []any {
	if len(args) == 0 {
		return args
	}
	out := make([]any, len(args))
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			out[i] = args[i]
			continue
		}
		key, ok := args[i].(string)
		if !ok {
			out[i] = args[i]
			out[i+1] = args[i+1]
			continue
		}
		out[i] = key
		out[i+1] = RedactValue(key, args[i+1])
	}
	return out
}

func newHandler(w io.Writer, level string, format string) slog.Handler {
	inner := baseHandler(w, level, format)
	return &RedactingHandler{inner: inner}
}
