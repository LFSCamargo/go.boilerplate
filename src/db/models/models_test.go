package models_test

import (
	"testing"

	"go.boilerplate/src/db/models"
)

func TestModels_TableNames(t *testing.T) {
	tests := []struct {
		model    any
		expected string
	}{
		{models.User{}, "users"},
		{models.UserRegistrationConfig{}, "user_registration_configs"},
		{models.RevokedToken{}, "revoked_tokens"},
		{models.OTP{}, "otps"},
	}

	for _, tt := range tests {
		namer, ok := tt.model.(interface{ TableName() string })
		if !ok {
			t.Fatalf("model %#v missing TableName", tt.model)
		}
		if got := namer.TableName(); got != tt.expected {
			t.Errorf("TableName() = %q, want %q", got, tt.expected)
		}
	}
}

func TestModels_AllIncludesEveryEntity(t *testing.T) {
	if len(models.All()) != 4 {
		t.Fatalf("expected 4 models, got %d", len(models.All()))
	}
}
