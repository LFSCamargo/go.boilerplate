package services

import (
	"errors"
	"testing"
)

func TestAuthErrorsAreDistinct(t *testing.T) {
	if errors.Is(ErrInvalidCredentials, ErrEmailTaken) {
		t.Fatal("errors should not wrap each other")
	}
	if ErrInvalidOTP.Error() == "" || ErrInvalidToken.Error() == "" {
		t.Fatal("errors need messages")
	}
}
