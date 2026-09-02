package services

import (
	"context"
	"strings"
	"time"

	"go.boilerplate/src/db/models"
	"go.boilerplate/src/mail"
)

func (s *AuthService) RecoverPassword(ctx context.Context, email string) error {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return err
	}

	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil || user == nil {
		return err
	}

	code, token, otp, err := s.newOTP(user.ID, models.OTPPurposePasswordReset, policy)
	if err != nil {
		return err
	}
	if err := s.otps.ConsumeActive(ctx, user.ID, models.OTPPurposePasswordReset); err != nil {
		return err
	}
	if err := s.otps.Create(ctx, otp); err != nil {
		return err
	}

	name := displayName(user)
	return s.mailer.SendPasswordReset(ctx, mail.PasswordResetParams{
		To:            user.Email,
		Name:          name,
		Code:          code,
		ExpiryMinutes: policy.OTPExpiryMinutes,
		ResetURL:      s.cfg.AppPublicURL + "/api/v1/auth/reset-password?token=" + token,
	})
}

func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return err
	}
	if len(newPassword) < policy.MinPasswordLength {
		return ErrWeakPassword
	}

	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidOTP
	}

	if _, err := s.consumeOTP(ctx, user, models.OTPPurposePasswordReset, code, policy); err != nil {
		return err
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.users.Update(ctx, user)
}

func (s *AuthService) ResetPasswordToken(ctx context.Context, token, newPassword string) error {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return err
	}
	if len(newPassword) < policy.MinPasswordLength {
		return ErrWeakPassword
	}
	otp, err := s.otps.FindActiveByTokenHash(ctx, HashSecret(strings.TrimSpace(token)))
	if err != nil {
		return err
	}
	if otp == nil || otp.Purpose != models.OTPPurposePasswordReset {
		return ErrInvalidToken
	}
	now := time.Now().UTC()
	otp.ConsumedAt = &now
	if err := s.otps.Save(ctx, otp); err != nil {
		return err
	}
	user, err := s.users.FindByID(ctx, otp.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.users.Update(ctx, user)
}
