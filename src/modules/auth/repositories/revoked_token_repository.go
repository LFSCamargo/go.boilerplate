package repositories

import (
	"context"
	"errors"
	"time"

	"go.boilerplate/src/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RevokedTokenRepository struct {
	db *gorm.DB
}

func NewRevokedTokenRepository(db *gorm.DB) *RevokedTokenRepository {
	return &RevokedTokenRepository{db: db}
}

func (r *RevokedTokenRepository) Revoke(ctx context.Context, userID uuid.UUID, jti string, expiresAt time.Time) error {
	row := models.RevokedToken{
		UserID:    userID,
		JTI:       jti,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *RevokedTokenRepository) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var row models.RevokedToken
	err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
