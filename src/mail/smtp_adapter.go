package mail

import (
	"context"

	smtpadapter "go.boilerplate/src/mail/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Security string
	User     string
	Password string
	From     string
}

type smtpSender struct {
	client *smtpadapter.Client
}

func newSMTPSender(cfg SMTPConfig) (*smtpSender, error) {
	client, err := smtpadapter.NewClient(smtpadapter.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Security: cfg.Security,
		User:     cfg.User,
		Password: cfg.Password,
		From:     cfg.From,
	})
	if err != nil {
		return nil, err
	}
	return &smtpSender{client: client}, nil
}

func (s *smtpSender) Send(ctx context.Context, msg Message) error {
	return s.client.Send(ctx, smtpadapter.Message{
		To:      msg.To,
		Subject: msg.Subject,
		HTML:    msg.HTML,
		Text:    msg.Text,
	})
}
