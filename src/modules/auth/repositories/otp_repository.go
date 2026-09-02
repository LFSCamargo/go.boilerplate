package repositories

import (
	"context"
	"errors"
	"time"

	"go.boilerplate/src/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *OTPRepository) ConsumeActive(ctx context.Context, userID uuid.UUID, purpose models.OTPPurpose) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&models.OTP{}).
		Where("user_id = ? AND purpose = ? AND consumed_at IS NULL", userID, purpose).
		Updates(map[string]any{"consumed_at": now}).Error
}

func (r *OTPRepository) FindLatestActive(ctx context.Context, userID uuid.UUID, purpose models.OTPPurpose) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ? AND consumed_at IS NULL AND expires_at > ?", userID, purpose, time.Now().UTC()).
		Order("created_at DESC").
		First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *OTPRepository) FindActiveByTokenHash(ctx context.Context, tokenHash string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND consumed_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *OTPRepository) Save(ctx context.Context, otp *models.OTP) error {
	return r.db.WithContext(ctx).Save(otp).Error
}
