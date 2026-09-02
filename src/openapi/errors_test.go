package openapi

import (
	"testing"

	"go.boilerplate/src/modules/auth/services"
)

func TestFromAuth_MapsKnownErrors(t *testing.T) {
	err := FromAuth(services.ErrEmailTaken)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("type %T", err)
	}
	if apiErr.GetStatus() != 409 || apiErr.Code != "email_taken" {
		t.Fatalf("%+v", apiErr)
	}
}
