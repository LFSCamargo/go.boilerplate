package models

// All returns every GORM model for registration, tooling, and tests.
func All() []any {
	return []any{
		&User{},
		&UserRegistrationConfig{},
		&RevokedToken{},
		&OTP{},
	}
}
