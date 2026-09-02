package openapi

import (
	"net/http"

	z "github.com/Oudwins/zog"
)

type validator interface {
	Validate(dataPtr any, options ...z.ExecOption) z.ZogIssueList
}

func Validate[T any](schema validator, dest T) error {
	if issues := schema.Validate(&dest); issues != nil {
		return Status(http.StatusBadRequest, "invalid_body")
	}
	return nil
}

func ValidateResponse[T any](schema validator, dest T) error {
	if issues := schema.Validate(&dest); issues != nil {
		return Status(http.StatusInternalServerError, "invalid_response")
	}
	return nil
}
