package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/mail"

	"github.com/google/uuid"
)

func (s *AuthService) issueVerification(ctx context.Context, user *models.User, policy *models.UserRegistrationConfig) error {
	code, token, otp, err := s.newOTP(user.ID, models.OTPPurposeEmailVerify, policy)
	if err != nil {
		return err
	}
	if err := s.otps.ConsumeActive(ctx, user.ID, models.OTPPurposeEmailVerify); err != nil {
		return err
	}
	if err := s.otps.Create(ctx, otp); err != nil {
		return err
	}

	return s.mailer.SendVerifyEmail(ctx, mail.VerifyEmailParams{
		To:            user.Email,
		Name:          displayName(user),
		Code:          code,
		ExpiryMinutes: policy.OTPExpiryMinutes,
		VerifyURL:     s.cfg.AppPublicURL + "/api/v1/auth/verify-email?token=" + token,
	})
}

func (s *AuthService) newOTP(userID uuid.UUID, purpose models.OTPPurpose, policy *models.UserRegistrationConfig) (code, token string, otp *models.OTP, err error) {
	code, err = GenerateOTPCode()
	if err != nil {
		return "", "", nil, err
	}
	token, err = GenerateMagicToken()
	if err != nil {
		return "", "", nil, err
	}
	expiry := policy.OTPExpiryMinutes
	if expiry <= 0 {
		expiry = 10
	}
	tokenHash := HashSecret(token)
	otp = &models.OTP{
		UserID:    userID,
		Purpose:   purpose,
		CodeHash:  HashSecret(code),
		TokenHash: &tokenHash,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiry) * time.Minute),
	}
	return code, token, otp, nil
}

func (s *AuthService) consumeOTP(ctx context.Context, user *models.User, purpose models.OTPPurpose, code string, policy *models.UserRegistrationConfig) (*models.OTP, error) {
	otp, err := s.otps.FindLatestActive(ctx, user.ID, purpose)
	if err != nil {
		return nil, err
	}
	if otp == nil {
		return nil, ErrInvalidOTP
	}

	maxAttempts := policy.MaxOTPAttempts
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	if otp.Attempts >= maxAttempts {
		return nil, ErrInvalidOTP
	}

	if otp.CodeHash != HashSecret(strings.TrimSpace(code)) {
		otp.Attempts++
		if saveErr := s.otps.Save(ctx, otp); saveErr != nil {
			return nil, saveErr
		}
		return nil, ErrInvalidOTP
	}

	now := time.Now().UTC()
	otp.ConsumedAt = &now
	if err := s.otps.Save(ctx, otp); err != nil {
		return nil, err
	}
	return otp, nil
}

func (s *AuthService) issueSession(user *models.User) (*TokenPair, error) {
	token, _, expiresAt, err := IssueAccessToken(s.cfg.BearerSecret, user.ID, user.Email, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("issue token: %w", err)
	}
	return &TokenPair{AccessToken: token, ExpiresAt: expiresAt}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func displayName(user *models.User) string {
	if user.DisplayName != nil && *user.DisplayName != "" {
		return *user.DisplayName
	}
	return "there"
}

func toAuthUser(user *models.User) *AuthUser {
	return &AuthUser{
		ID:            user.ID,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		AvatarURL:     user.AvatarURL,
		EmailVerified: user.EmailVerified,
	}
}
