package authHandlers

import (
	"context"
	"mime/multipart"
	"time"

	"go.boilerplate/src/db/models"

	"github.com/google/uuid"
)

type memUsers struct {
	users map[string]*models.User
}

func (m *memUsers) Create(_ context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	clone := *user
	if m.users == nil {
		m.users = map[string]*models.User{}
	}
	m.users[user.Email] = &clone
	return nil
}

func (m *memUsers) FindByEmail(_ context.Context, email string) (*models.User, error) {
	if u, ok := m.users[email]; ok {
		clone := *u
		return &clone, nil
	}
	return nil, nil
}

func (m *memUsers) FindByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			clone := *u
			return &clone, nil
		}
	}
	return nil, nil
}

func (m *memUsers) Update(_ context.Context, user *models.User) error {
	clone := *user
	m.users[user.Email] = &clone
	return nil
}

func (m *memUsers) RegistrationConfig(context.Context) (*models.UserRegistrationConfig, error) {
	return &models.UserRegistrationConfig{
		RequireEmailVerification: true,
		MinPasswordLength:        8,
		AllowRegistration:        true,
		OTPExpiryMinutes:         10,
		MaxOTPAttempts:           5,
	}, nil
}

func (m *memUsers) SaveAvatar(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}
	return "http://localhost:8080/uploads/avatars/" + file.Filename, nil
}

type memOTPs struct{ rows []*models.OTP }

func (m *memOTPs) Create(_ context.Context, otp *models.OTP) error {
	if otp.ID == uuid.Nil {
		otp.ID = uuid.New()
	}
	clone := *otp
	m.rows = append(m.rows, &clone)
	return nil
}

func (m *memOTPs) ConsumeActive(_ context.Context, userID uuid.UUID, purpose models.OTPPurpose) error {
	now := time.Now().UTC()
	for _, row := range m.rows {
		if row.UserID == userID && row.Purpose == purpose && row.ConsumedAt == nil {
			row.ConsumedAt = &now
		}
	}
	return nil
}

func (m *memOTPs) FindLatestActive(_ context.Context, userID uuid.UUID, purpose models.OTPPurpose) (*models.OTP, error) {
	now := time.Now().UTC()
	var found *models.OTP
	for _, row := range m.rows {
		if row.UserID == userID && row.Purpose == purpose && row.ConsumedAt == nil && row.ExpiresAt.After(now) {
			found = row
		}
	}
	if found == nil {
		return nil, nil
	}
	clone := *found
	return &clone, nil
}

func (m *memOTPs) FindActiveByTokenHash(_ context.Context, tokenHash string) (*models.OTP, error) {
	now := time.Now().UTC()
	for _, row := range m.rows {
		if row.TokenHash != nil && *row.TokenHash == tokenHash && row.ConsumedAt == nil && row.ExpiresAt.After(now) {
			clone := *row
			return &clone, nil
		}
	}
	return nil, nil
}

func (m *memOTPs) Save(_ context.Context, otp *models.OTP) error {
	for i, row := range m.rows {
		if row.ID == otp.ID {
			clone := *otp
			m.rows[i] = &clone
		}
	}
	return nil
}

type memRevoked struct{}

func (memRevoked) Revoke(context.Context, uuid.UUID, string, time.Time) error { return nil }
func (memRevoked) IsRevoked(context.Context, string) (bool, error)            { return false, nil }
