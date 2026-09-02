package middleware

import (
	z "github.com/Oudwins/zog"
	"github.com/gofiber/fiber/v3"
)

// Validator is a Zog schema that can validate an already-built value.
type Validator interface {
	Validate(dataPtr any, options ...z.ExecOption) z.ZogIssueList
}

// WriteJSON validates dest with a Zog schema, then writes it as JSON.
func WriteJSON[T any](c fiber.Ctx, schema Validator, dest T) error {
	return WriteJSONStatus(c, fiber.StatusOK, schema, dest)
}

// WriteJSONStatus is WriteJSON with an explicit status code.
func WriteJSONStatus[T any](c fiber.Ctx, status int, schema Validator, dest T) error {
	if issues := schema.Validate(&dest); issues != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "invalid_response",
			"issues": z.Issues.Flatten(issues),
		})
	}
	return c.Status(status).JSON(dest)
}
