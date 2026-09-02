package requests

import "testing"

func TestParseJSON_LoginRejectsInvalidEmail(t *testing.T) {
	var body Login
	issues := ParseJSON([]byte(`{"email":"nope","password":"x"}`), LoginSchema, &body)
	if issues == nil {
		t.Fatal("expected validation issues")
	}
	flat := Flatten(issues)
	if _, ok := flat["email"]; !ok {
		if _, ok := flat["Email"]; !ok {
			t.Fatalf("expected email issue, got %#v", flat)
		}
	}
}

func TestParseJSON_RegisterAcceptsValidBody(t *testing.T) {
	var body Register
	issues := ParseJSON([]byte(`{"email":"ada@example.com","password":"password12","display_name":"Ada"}`), RegisterSchema, &body)
	if issues != nil {
		t.Fatalf("unexpected issues: %#v", Flatten(issues))
	}
	if body.Email != "ada@example.com" || body.Password != "password12" {
		t.Fatalf("parsed %#v", body)
	}
}
