package repositories

import (
	"context"
	"errors"
	"strings"

	"go.boilerplate/src/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db         *gorm.DB
	uploadsDir string
	publicURL  string
}

func NewUserRepository(db *gorm.DB, uploadsDir, publicURL string) *UserRepository {
	return &UserRepository{db: db, uploadsDir: uploadsDir, publicURL: publicURL}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(email))).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) RegistrationConfig(ctx context.Context) (*models.UserRegistrationConfig, error) {
	var cfg models.UserRegistrationConfig
	err := r.db.WithContext(ctx).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.UserRegistrationConfig{
			RequireEmailVerification: true,
			MinPasswordLength:        8,
			AllowRegistration:        true,
			OTPExpiryMinutes:         10,
			MaxOTPAttempts:           5,
		}, nil
	}
	return &cfg, err
}
