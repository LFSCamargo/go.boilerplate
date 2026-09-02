package requests

import "testing"

func TestRegisterSchema_RequiresMinPassword(t *testing.T) {
	var body Register
	if issues := ParseJSON([]byte(`{"email":"ada@example.com","password":"short"}`), RegisterSchema, &body); issues == nil {
		t.Fatal("expected issues")
	}
}
