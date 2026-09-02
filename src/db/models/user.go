package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string    `gorm:"size:255;not null"`
	PasswordHash  string    `gorm:"size:255;not null;column:password_hash"`
	DisplayName   *string   `gorm:"size:100;column:display_name"`
	AvatarURL     *string   `gorm:"size:512;column:avatar_url"`
	EmailVerified bool      `gorm:"not null;default:false;column:email_verified"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime"`
}

func (User) TableName() string { return "users" }
