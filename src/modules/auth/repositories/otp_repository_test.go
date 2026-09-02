package repositories

import "testing"

func TestNewOTPRepository(t *testing.T) {
	if NewOTPRepository(nil) == nil {
		t.Fatal("expected repository")
	}
}
