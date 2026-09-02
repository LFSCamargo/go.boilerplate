package services

import (
	"context"
	"mime/multipart"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/mail"

	"github.com/google/uuid"
)

type UserStore interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	RegistrationConfig(ctx context.Context) (*models.UserRegistrationConfig, error)
	SaveAvatar(file *multipart.FileHeader) (string, error)
}

type OTPStore interface {
	Create(ctx context.Context, otp *models.OTP) error
	ConsumeActive(ctx context.Context, userID uuid.UUID, purpose models.OTPPurpose) error
	FindLatestActive(ctx context.Context, userID uuid.UUID, purpose models.OTPPurpose) (*models.OTP, error)
	FindActiveByTokenHash(ctx context.Context, tokenHash string) (*models.OTP, error)
	Save(ctx context.Context, otp *models.OTP) error
}

type RevokedStore interface {
	Revoke(ctx context.Context, userID uuid.UUID, jti string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

type Mailer interface {
	SendVerifyEmail(ctx context.Context, p mail.VerifyEmailParams) error
	SendPasswordReset(ctx context.Context, p mail.PasswordResetParams) error
}
