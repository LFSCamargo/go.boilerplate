package middleware

import (
	"bytes"
	"reflect"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
	"github.com/gofiber/fiber/v3"
)

const localPayload = "validated_payload"

// Schema is any Zog schema that can parse into a dest pointer.
type Schema interface {
	Parse(data any, destPtr any, options ...z.ExecOption) z.ZogIssueList
}

// ValidateBody parses JSON or form bodies with a Zog schema and injects the
// typed value for Payload[T].
func ValidateBody[T any](schema Schema) fiber.Handler {
	return func(c fiber.Ctx) error {
		var dest T
		var issues z.ZogIssueList

		raw := bytes.TrimSpace(c.Body())
		switch {
		case isForm(c.Get("Content-Type")):
			issues = schema.Parse(valuesFromTags[T](func(key string) string {
				return c.FormValue(key)
			}), &dest)
		case len(raw) == 0:
			issues = schema.Parse(map[string]any{}, &dest)
		default:
			issues = schema.Parse(zjson.Decode(bytes.NewReader(raw)), &dest)
		}
		if issues != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "invalid_body",
				"issues": z.Issues.Flatten(issues),
			})
		}

		c.Locals(localPayload, dest)
		return c.Next()
	}
}

// ValidateQuery parses query parameters with a Zog schema and injects Payload[T].
func ValidateQuery[T any](schema Schema) fiber.Handler {
	return func(c fiber.Ctx) error {
		var dest T
		if issues := schema.Parse(valuesFromTags[T](func(key string) string {
			return c.Query(key)
		}), &dest); issues != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "invalid_query",
				"issues": z.Issues.Flatten(issues),
			})
		}
		c.Locals(localPayload, dest)
		return c.Next()
	}
}

// Payload returns the value parsed by ValidateBody or ValidateQuery.
func Payload[T any](c fiber.Ctx) T {
	value, _ := c.Locals(localPayload).(T)
	return value
}

func isForm(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.HasPrefix(ct, "multipart/form-data") ||
		strings.HasPrefix(ct, "application/x-www-form-urlencoded")
}

func valuesFromTags[T any](get func(string) string) map[string]any {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	data := make(map[string]any, typ.NumField())
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		key := firstTag(field, "form", "query", "json")
		if key == "" || key == "-" {
			key = field.Name
		}
		data[field.Name] = get(key)
	}
	return data
}

func firstTag(field reflect.StructField, names ...string) string {
	for _, name := range names {
		if tag := field.Tag.Get(name); tag != "" {
			return strings.Split(tag, ",")[0]
		}
	}
	return ""
}
