package mail

import (
	"context"
	"fmt"

	"go.boilerplate/src/config"
	"go.boilerplate/src/mail/reactemail"
)

type Service struct {
	sender   Sender
	renderer Renderer
	from     string
}

func NewService(sender Sender, renderer Renderer, from string) *Service {
	return &Service{
		sender:   sender,
		renderer: renderer,
		from:     from,
	}
}

func NewFromConfig(cfg *config.Config) (*Service, error) {
	return NewFromSMTPConfig(SMTPConfigFromApp(cfg), cfg.EmailsDir, cfg.EmailAssetsBaseURL)
}

func NewFromSMTPConfig(smtpCfg SMTPConfig, emailsDir, assetsBaseURL string) (*Service, error) {
	sender, err := newSMTPSender(smtpCfg)
	if err != nil {
		return nil, err
	}

	renderer, err := reactemail.NewRenderer(emailsDir, assetsBaseURL)
	if err != nil {
		return nil, err
	}

	return NewService(sender, renderer, smtpCfg.From), nil
}

func (s *Service) SendVerifyEmail(ctx context.Context, p VerifyEmailParams) error {
	if p.ExpiryMinutes == 0 {
		p.ExpiryMinutes = 10
	}

	html, err := s.renderer.Render(TemplateVerifyEmail, map[string]any{
		"name":          p.Name,
		"code":          p.Code,
		"expiryMinutes": p.ExpiryMinutes,
		"companyName":   BrandName,
		"verifyUrl":     p.VerifyURL,
	})
	if err != nil {
		return fmt.Errorf("render verify_email: %w", err)
	}

	return s.sender.Send(ctx, Message{
		To:      p.To,
		Subject: "Verify your Go Boilerplate email",
		HTML:    html,
	})
}

func (s *Service) SendPasswordReset(ctx context.Context, p PasswordResetParams) error {
	if p.ExpiryMinutes == 0 {
		p.ExpiryMinutes = 10
	}

	html, err := s.renderer.Render(TemplatePasswordReset, map[string]any{
		"name":          p.Name,
		"code":          p.Code,
		"expiryMinutes": p.ExpiryMinutes,
		"companyName":   BrandName,
		"resetUrl":      p.ResetURL,
	})
	if err != nil {
		return fmt.Errorf("render password_reset: %w", err)
	}

	return s.sender.Send(ctx, Message{
		To:      p.To,
		Subject: "Reset your Go Boilerplate password",
		HTML:    html,
	})
}

func (s *Service) Render(template string, props map[string]any) (string, error) {
	return s.renderer.Render(template, props)
}

func (s *Service) Send(ctx context.Context, msg Message) error {
	return s.sender.Send(ctx, msg)
}
