package log

import (
	"context"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
)

const redacted = "[REDACTED]"

var (
	jwtPattern      = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
	bearerPattern   = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)
	postgresPattern = regexp.MustCompile(`postgres(?:ql)?://[^:\s]+:[^@\s]+@`)
)

var redactKeys = map[string]struct{}{
	"password":            {},
	"new_password":        {},
	"old_password":        {},
	"token":               {},
	"access_token":        {},
	"refresh_token":       {},
	"authorization":       {},
	"cookie":              {},
	"set-cookie":          {},
	"code":                {},
	"otp":                 {},
	"bearer_secret":       {},
	"secret":              {},
	"smtp_password":       {},
	"smtp_user":           {},
	"picture":             {},
}

// RedactingHandler wraps a slog.Handler and scrubs sensitive attributes.
type RedactingHandler struct {
	inner slog.Handler
}

func (h *RedactingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *RedactingHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.NumAttrs() == 0 {
		return h.inner.Handle(ctx, r)
	}
	rec := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.Attrs(func(a slog.Attr) bool {
		rec.AddAttrs(redactAttr(a))
		return true
	})
	return h.inner.Handle(ctx, rec)
}

func (h *RedactingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	redacted := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		redacted[i] = redactAttr(a)
	}
	return &RedactingHandler{inner: h.inner.WithAttrs(redacted)}
}

func (h *RedactingHandler) WithGroup(name string) slog.Handler {
	return &RedactingHandler{inner: h.inner.WithGroup(name)}
}

func redactAttr(a slog.Attr) slog.Attr {
	if a.Equal(slog.Attr{}) {
		return a
	}
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		for i, ga := range attrs {
			attrs[i] = redactAttr(ga)
		}
		return slog.Group(a.Key, groupArgs(attrs)...)
	}
	return slog.Any(a.Key, RedactValue(a.Key, a.Value.Any()))
}

func groupArgs(attrs []slog.Attr) []any {
	args := make([]any, 0, len(attrs)*2)
	for _, a := range attrs {
		args = append(args, a.Key, a.Value.Any())
	}
	return args
}

// RedactValue applies field-aware and pattern-based redaction.
func RedactValue(key string, value any) any {
	switch v := value.(type) {
	case string:
		return redactString(key, v)
	case []byte:
		return redacted
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, val := range v {
			out[k] = RedactValue(k, val)
		}
		return out
	case map[string]string:
		out := make(map[string]string, len(v))
		for k, val := range v {
			out[k] = RedactValue(k, val).(string)
		}
		return out
	default:
		if shouldRedactKey(key) {
			return redacted
		}
		return value
	}
}

func redactString(key, value string) string {
	if value == "" {
		return value
	}
	if strings.EqualFold(key, "email") {
		return maskEmail(value)
	}
	if shouldRedactKey(key) {
		if strings.EqualFold(key, "authorization") {
			return redactStringPatterns(value)
		}
		return redacted
	}
	return redactStringPatterns(value)
}

func shouldRedactKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	if _, ok := redactKeys[normalized]; ok {
		return true
	}
	return strings.Contains(normalized, "password") ||
		strings.Contains(normalized, "token") ||
		strings.Contains(normalized, "secret") ||
		strings.Contains(normalized, "authorization")
}

func redactStringPatterns(value string) string {
	value = bearerPattern.ReplaceAllString(value, "Bearer "+redacted)
	value = jwtPattern.ReplaceAllString(value, redacted)
	value = postgresPattern.ReplaceAllString(value, "postgres://[REDACTED]:[REDACTED]@")
	return value
}

func maskEmail(email string) string {
	email = strings.TrimSpace(email)
	at := strings.LastIndex(email, "@")
	if at <= 0 {
		return redacted
	}
	local := email[:at]
	domain := email[at:]
	if len(local) == 1 {
		return local + "***" + domain
	}
	return local[:1] + "***" + domain
}

// RedactQuery scrubs sensitive query parameters from a URL or raw query string.
func RedactQuery(raw string) string {
	if raw == "" {
		return raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return redactStringPatterns(raw)
	}
	q := parsed.Query()
	for key := range q {
		if shouldRedactKey(key) {
			q.Set(key, redacted)
		}
	}
	parsed.RawQuery = q.Encode()
	return parsed.String()
}
