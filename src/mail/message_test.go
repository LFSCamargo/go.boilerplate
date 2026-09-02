package mail

import "testing"

func TestBrandAndTemplateNames(t *testing.T) {
	if BrandName != "Go Boilerplate" {
		t.Fatalf("got %s", BrandName)
	}
	if TemplateVerifyEmail != "verify_email" || TemplatePasswordReset != "password_reset" {
		t.Fatal("unexpected template slugs")
	}
}
