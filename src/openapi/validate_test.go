package openapi

import (
	"testing"

	z "github.com/Oudwins/zog"
)

func TestValidate_RejectsInvalid(t *testing.T) {
	schema := z.Struct(z.Shape{"Email": z.String().Required().Email()})
	type body struct {
		Email string
	}
	if err := Validate(schema, body{Email: "nope"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_AcceptsValid(t *testing.T) {
	schema := z.Struct(z.Shape{"Email": z.String().Required().Email()})
	type body struct {
		Email string
	}
	if err := Validate(schema, body{Email: "ada@example.com"}); err != nil {
		t.Fatal(err)
	}
}
