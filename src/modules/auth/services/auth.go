package services

import (
	"context"
	"mime/multipart"
	"strings"
	"time"

	"go.boilerplate/src/db/models"

	"github.com/google/uuid"
)

type Config struct {
	BearerSecret   string
	AppPublicURL   string
	AccessTokenTTL time.Duration
}

type AuthService struct {
	users   UserStore
	otps    OTPStore
	revoked RevokedStore
	mailer  Mailer
	cfg     Config
}

func NewAuthService(users UserStore, otps OTPStore, revoked RevokedStore, mailer Mailer, cfg Config) *AuthService {
	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = 24 * time.Hour
	}
	cfg.AppPublicURL = strings.TrimRight(cfg.AppPublicURL, "/")
	return &AuthService{users: users, otps: otps, revoked: revoked, mailer: mailer, cfg: cfg}
}

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
	Picture     *multipart.FileHeader
}

type AuthUser struct {
	ID            uuid.UUID
	Email         string
	DisplayName   *string
	AvatarURL     *string
	EmailVerified bool
}

type TokenPair struct {
	AccessToken string
	ExpiresAt   time.Time
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthUser, error) {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !policy.AllowRegistration {
		return nil, ErrRegistrationClosed
	}

	email := normalizeEmail(in.Email)
	if email == "" || in.Password == "" {
		return nil, ErrInvalidCredentials
	}
	if len(in.Password) < policy.MinPasswordLength {
		return nil, ErrWeakPassword
	}

	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hash,
	}
	if name := strings.TrimSpace(in.DisplayName); name != "" {
		user.DisplayName = &name
	}
	if in.Picture != nil {
		url, err := s.users.SaveAvatar(in.Picture)
		if err != nil {
			return nil, ErrInvalidPicture
		}
		if url != "" {
			user.AvatarURL = &url
		}
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	if err := s.issueVerification(ctx, user, policy); err != nil {
		return nil, err
	}

	return toAuthUser(user), nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthUser, *TokenPair, error) {
	policy, err := s.users.RegistrationConfig(ctx)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, nil, err
	}
	if user == nil || CheckPassword(user.PasswordHash, password) != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if policy.RequireEmailVerification && !user.EmailVerified {
		_ = s.issueVerification(ctx, user, policy)
	}

	token, err := s.issueSession(user)
	if err != nil {
		return nil, nil, err
	}
	return toAuthUser(user), token, nil
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (*AuthUser, *Claims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, nil, ErrInvalidToken
	}

	claims, err := ParseAccessToken(s.cfg.BearerSecret, token)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	revoked, err := s.revoked.IsRevoked(ctx, claims.ID)
	if err != nil {
		return nil, nil, err
	}
	if revoked {
		return nil, nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrInvalidToken
	}

	return toAuthUser(user), claims, nil
}

func (s *AuthService) Logout(ctx context.Context, bearer string) error {
	claims, err := ParseAccessToken(s.cfg.BearerSecret, bearer)
	if err != nil {
		return ErrInvalidToken
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ErrInvalidToken
	}
	if claims.ExpiresAt == nil {
		return ErrInvalidToken
	}
	return s.revoked.Revoke(ctx, userID, claims.ID, claims.ExpiresAt.Time)
}
