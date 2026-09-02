package models

import "testing"

func TestOTPTableNameAndPurposes(t *testing.T) {
	if got := (OTP{}).TableName(); got != "otps" {
		t.Fatalf("got %s", got)
	}
	if OTPPurposeEmailVerify != "email_verify" || OTPPurposePasswordReset != "password_reset" {
		t.Fatal("unexpected OTP purposes")
	}
}
