package requests

import "testing"

func TestLoginSchema_RequiresPassword(t *testing.T) {
	var body Login
	if issues := ParseJSON([]byte(`{"email":"ada@example.com"}`), LoginSchema, &body); issues == nil {
		t.Fatal("expected issues")
	}
}
