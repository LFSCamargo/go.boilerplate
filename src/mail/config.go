package mail

import "go.boilerplate/src/config"

func SMTPConfigFromApp(cfg *config.Config) SMTPConfig {
	return SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Security: cfg.SMTPSecurity,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.MailFrom,
	}
}

func IsConfigured(cfg *config.Config) bool {
	return cfg.SMTPHost != "" && cfg.MailFrom != ""
}
