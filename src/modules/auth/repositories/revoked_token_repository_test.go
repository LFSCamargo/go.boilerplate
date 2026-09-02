package repositories

import "testing"

func TestNewRevokedTokenRepository(t *testing.T) {
	if NewRevokedTokenRepository(nil) == nil {
		t.Fatal("expected repository")
	}
}
