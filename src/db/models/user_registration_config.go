package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRegistrationConfig holds application-wide registration policy (singleton row).
type UserRegistrationConfig struct {
	ID                       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequireEmailVerification bool      `gorm:"not null;default:true;column:require_email_verification"`
	MinPasswordLength        int       `gorm:"not null;default:8;column:min_password_length"`
	AllowRegistration        bool      `gorm:"not null;default:true;column:allow_registration"`
	OTPExpiryMinutes         int       `gorm:"not null;default:10;column:otp_expiry_minutes"`
	MaxOTPAttempts           int       `gorm:"not null;default:5;column:max_otp_attempts"`
	CreatedAt                time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt                time.Time `gorm:"not null;autoUpdateTime"`
}

func (UserRegistrationConfig) TableName() string { return "user_registration_configs" }
