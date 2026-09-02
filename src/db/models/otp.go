package models

import (
	"time"

	"github.com/google/uuid"
)

type OTPPurpose string

const (
	OTPPurposeEmailVerify   OTPPurpose = "email_verify"
	OTPPurposePasswordReset OTPPurpose = "password_reset"
)

type OTP struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_otps_user_purpose;column:user_id"`
	User       User       `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	Purpose    OTPPurpose `gorm:"type:otp_purpose;not null;index:idx_otps_user_purpose"`
	CodeHash   string     `gorm:"size:255;not null;column:code_hash"`
	TokenHash  *string    `gorm:"size:255;column:token_hash"`
	ExpiresAt  time.Time  `gorm:"not null;index;column:expires_at"`
	ConsumedAt *time.Time `gorm:"column:consumed_at"`
	Attempts   int        `gorm:"not null;default:0"`
	CreatedAt  time.Time  `gorm:"not null;autoCreateTime"`
}

func (OTP) TableName() string { return "otps" }
