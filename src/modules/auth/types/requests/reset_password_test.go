package requests

import "testing"

func TestResetPassword_HasTokenAndCode(t *testing.T) {
	token := ResetPassword{Token: "abc", NewPassword: "password12"}
	if !token.HasToken() || token.HasCode() {
		t.Fatalf("token body: %+v", token)
	}
	code := ResetPassword{Email: "ada@example.com", Code: "123456", NewPassword: "password12"}
	if !code.HasCode() || code.HasToken() {
		t.Fatalf("code body: %+v", code)
	}
}
