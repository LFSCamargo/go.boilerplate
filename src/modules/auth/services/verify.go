package services

import (
	"context"
	"strings"
	"time"

	"go.boilerplate/src/db/models"
)

func (s *AuthService) ResendVerification(ctx context.Context, email string) error {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return err
	}
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil || user == nil {
		return err
	}
	if user.EmailVerified {
		return ErrAlreadyVerified
	}
	return s.issueVerification(ctx, user, policy)
}

func (s *AuthService) VerifyEmailCode(ctx context.Context, email, code string) error {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return err
	}
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidOTP
	}
	if user.EmailVerified {
		return ErrAlreadyVerified
	}
	if _, err := s.consumeOTP(ctx, user, models.OTPPurposeEmailVerify, code, policy); err != nil {
		return err
	}
	user.EmailVerified = true
	return s.users.Update(ctx, user)
}

func (s *AuthService) VerifyEmailToken(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidToken
	}
	otp, err := s.otps.FindActiveByTokenHash(ctx, HashSecret(token))
	if err != nil {
		return err
	}
	if otp == nil || otp.Purpose != models.OTPPurposeEmailVerify {
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
	user.EmailVerified = true
	return s.users.Update(ctx, user)
}
