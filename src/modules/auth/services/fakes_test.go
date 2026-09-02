package services

import (
	"context"
	"mime/multipart"
	"sync"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/mail"

	"github.com/google/uuid"
)

type fakeUsers struct {
	mu     sync.Mutex
	byID   map[uuid.UUID]*models.User
	policy *models.UserRegistrationConfig
}

func newFakeUsers() *fakeUsers {
	return &fakeUsers{
		byID: make(map[uuid.UUID]*models.User),
		policy: &models.UserRegistrationConfig{
			RequireEmailVerification: true,
			MinPasswordLength:        8,
			AllowRegistration:        true,
			OTPExpiryMinutes:         10,
			MaxOTPAttempts:           5,
		},
	}
}

func (f *fakeUsers) Create(_ context.Context, user *models.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	clone := *user
	f.byID[user.ID] = &clone
	return nil
}

func (f *fakeUsers) FindByEmail(_ context.Context, email string) (*models.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byID {
		if u.Email == normalizeEmail(email) {
			clone := *u
			return &clone, nil
		}
	}
	return nil, nil
}

func (f *fakeUsers) FindByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.byID[id]
	if !ok {
		return nil, nil
	}
	clone := *u
	return &clone, nil
}

func (f *fakeUsers) Update(_ context.Context, user *models.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	clone := *user
	f.byID[user.ID] = &clone
	return nil
}

func (f *fakeUsers) RegistrationConfig(context.Context) (*models.UserRegistrationConfig, error) {
	return f.policy, nil
}

func (f *fakeUsers) SaveAvatar(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}
	return "http://localhost:8080/uploads/avatars/" + file.Filename, nil
}

type fakeOTPs struct {
	mu   sync.Mutex
	rows []*models.OTP
}

func (f *fakeOTPs) Create(_ context.Context, otp *models.OTP) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if otp.ID == uuid.Nil {
		otp.ID = uuid.New()
	}
	clone := *otp
	f.rows = append(f.rows, &clone)
	return nil
}

func (f *fakeOTPs) ConsumeActive(_ context.Context, userID uuid.UUID, purpose models.OTPPurpose) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now().UTC()
	for _, row := range f.rows {
		if row.UserID == userID && row.Purpose == purpose && row.ConsumedAt == nil {
			row.ConsumedAt = &now
		}
	}
	return nil
}

func (f *fakeOTPs) FindLatestActive(_ context.Context, userID uuid.UUID, purpose models.OTPPurpose) (*models.OTP, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now().UTC()
	var found *models.OTP
	for _, row := range f.rows {
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

func (f *fakeOTPs) FindActiveByTokenHash(_ context.Context, tokenHash string) (*models.OTP, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now().UTC()
	for _, row := range f.rows {
		if row.TokenHash != nil && *row.TokenHash == tokenHash && row.ConsumedAt == nil && row.ExpiresAt.After(now) {
			clone := *row
			return &clone, nil
		}
	}
	return nil, nil
}

func (f *fakeOTPs) Save(_ context.Context, otp *models.OTP) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, row := range f.rows {
		if row.ID == otp.ID {
			clone := *otp
			f.rows[i] = &clone
			return nil
		}
	}
	return nil
}

type captureMailer struct {
	mu     sync.Mutex
	verify []mail.VerifyEmailParams
	reset  []mail.PasswordResetParams
}

func (c *captureMailer) SendVerifyEmail(_ context.Context, p mail.VerifyEmailParams) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.verify = append(c.verify, p)
	return nil
}

func (c *captureMailer) SendPasswordReset(_ context.Context, p mail.PasswordResetParams) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reset = append(c.reset, p)
	return nil
}

type fakeRevoked struct {
	mu   sync.Mutex
	jtis map[string]bool
}

func newFakeRevoked() *fakeRevoked {
	return &fakeRevoked{jtis: map[string]bool{}}
}

func (f *fakeRevoked) Revoke(_ context.Context, _ uuid.UUID, jti string, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.jtis[jti] = true
	return nil
}

func (f *fakeRevoked) IsRevoked(_ context.Context, jti string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.jtis[jti], nil
}
