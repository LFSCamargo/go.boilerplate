package mail_test

import (
	"testing"

	"go.boilerplate/src/config"
	"go.boilerplate/src/mail"
)

func TestIsConfigured(t *testing.T) {
	cfg := &config.Config{
		SMTPHost: "localhost",
		MailFrom: "noreply@go.boilerplate",
	}
	if !mail.IsConfigured(cfg) {
		t.Fatal("expected mail to be configured")
	}

	cfg.SMTPHost = ""
	if mail.IsConfigured(cfg) {
		t.Fatal("expected mail to be disabled without smtp host")
	}
}

func TestSMTPConfigFromApp(t *testing.T) {
	cfg := &config.Config{
		SMTPHost:     "localhost",
		SMTPPort:     "1025",
		SMTPSecurity: "none",
		MailFrom:     "noreply@go.boilerplate",
	}

	smtpCfg := mail.SMTPConfigFromApp(cfg)
	if smtpCfg.Host != "localhost" || smtpCfg.From != "noreply@go.boilerplate" {
		t.Fatalf("unexpected smtp config: %+v", smtpCfg)
	}
}
