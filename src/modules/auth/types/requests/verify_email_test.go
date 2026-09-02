package requests

import "testing"

func TestVerifyEmailCodeSchema_RequiresSixDigits(t *testing.T) {
	var body VerifyEmailCode
	if issues := ParseJSON([]byte(`{"code":"12"}`), VerifyEmailCodeSchema, &body); issues == nil {
		t.Fatal("expected issues")
	}
	if issues := ParseJSON([]byte(`{"code":"123456"}`), VerifyEmailCodeSchema, &body); issues != nil {
		t.Fatalf("expected valid code: %v", issues)
	}
}

func TestVerifyEmailTokenSchema_RequiresToken(t *testing.T) {
	var body VerifyEmailToken
	if issues := ParseMap(map[string]any{"Token": ""}, VerifyEmailTokenSchema, &body); issues == nil {
		t.Fatal("expected issues")
	}
}
