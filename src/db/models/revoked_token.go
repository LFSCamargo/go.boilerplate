package models

import (
	"time"

	"github.com/google/uuid"
)

type RevokedToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;column:user_id"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	JTI       string    `gorm:"size:255;not null;uniqueIndex;column:jti"`
	ExpiresAt time.Time `gorm:"not null;index;column:expires_at"`
	RevokedAt time.Time `gorm:"not null;autoCreateTime;column:revoked_at"`
}

func (RevokedToken) TableName() string { return "revoked_tokens" }
