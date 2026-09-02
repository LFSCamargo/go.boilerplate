package mail

import "context"

const (
	TemplateVerifyEmail   = "verify_email"
	TemplatePasswordReset = "password_reset"
	BrandName             = "Go Boilerplate"
)

type Message struct {
	To      string
	Subject string
	HTML    string
	Text    string
}

type Renderer interface {
	Render(template string, props map[string]any) (html string, err error)
}

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

type Mailer interface {
	Renderer
	Sender
}

type VerifyEmailParams struct {
	To            string
	Name          string
	Code          string
	ExpiryMinutes int
	VerifyURL     string
}

type PasswordResetParams struct {
	To            string
	Name          string
	Code          string
	ExpiryMinutes int
	ResetURL      string
}
