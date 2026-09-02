package authRouter

import (
	"context"
	"mime/multipart"
	"testing"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/modules/auth/services"

	"github.com/google/uuid"
)

// in-memory stores duplicated for router wiring (keeps services test fakes unexported)

type memUsers struct {
	users map[string]*models.User
}

func (m *memUsers) Create(_ context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	clone := *user
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

type memRevoked struct {
	jtis map[string]bool
}

func (m *memRevoked) Revoke(_ context.Context, _ uuid.UUID, jti string, _ time.Time) error {
	if m.jtis == nil {
		m.jtis = map[string]bool{}
	}
	m.jtis[jti] = true
	return nil
}

func (m *memRevoked) IsRevoked(_ context.Context, jti string) (bool, error) {
	return m.jtis[jti], nil
}

func (m *memUsers) markVerified(email string) {
	if u, ok := m.users[email]; ok {
		u.EmailVerified = true
	}
}

func newRouterService(t *testing.T) *services.AuthService {
	t.Helper()
	svc, _ := newRouterBundle(t)
	return svc
}

func newRouterBundle(t *testing.T) (*services.AuthService, *memUsers) {
	t.Helper()
	users := &memUsers{users: map[string]*models.User{}}
	svc := services.NewAuthService(
		users,
		&memOTPs{},
		&memRevoked{jtis: map[string]bool{}},
		services.NoopMailer{},
		services.Config{
			BearerSecret: "dev-secret-key-at-least-32-characters-long",
			AppPublicURL: "http://localhost:8080",
		},
	)
	return svc, users
}
