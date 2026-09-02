package services

import "testing"

func TestHashPassword_RoundTrip(t *testing.T) {
	hash, err := HashPassword("s3cret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(hash, "s3cret-pass"); err != nil {
		t.Fatalf("expected password to match: %v", err)
	}
	if err := CheckPassword(hash, "wrong"); err == nil {
		t.Fatal("expected mismatch")
	}
}

func TestGenerateOTPCode_SixDigits(t *testing.T) {
	code, err := GenerateOTPCode()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Fatalf("got %q", code)
	}
}

func TestHashSecret_Stable(t *testing.T) {
	if HashSecret("abc") != HashSecret("abc") {
		t.Fatal("hash should be deterministic")
	}
	if HashSecret("abc") == HashSecret("abd") {
		t.Fatal("different inputs should not collide")
	}
}
