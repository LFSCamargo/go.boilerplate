package mail_test

import (
	"context"
	"testing"

	"go.boilerplate/src/mail"
)

type stubRenderer struct {
	lastTemplate string
	lastProps    map[string]any
	html         string
}

func (s *stubRenderer) Render(template string, props map[string]any) (string, error) {
	s.lastTemplate = template
	s.lastProps = props
	return s.html, nil
}

type stubSender struct {
	last mail.Message
}

func (s *stubSender) Send(_ context.Context, msg mail.Message) error {
	s.last = msg
	return nil
}

func TestSendVerifyEmail_UsesGoBoilerplateBrand(t *testing.T) {
	renderer := &stubRenderer{html: "<html>ok</html>"}
	sender := &stubSender{}
	svc := mail.NewService(sender, renderer, "noreply@go.boilerplate")

	err := svc.SendVerifyEmail(context.Background(), mail.VerifyEmailParams{
		To:        "ada@example.com",
		Name:      "Ada",
		Code:      "123456",
		VerifyURL: "http://localhost:8080/api/v1/auth/verify-email?token=abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if renderer.lastTemplate != mail.TemplateVerifyEmail {
		t.Fatalf("template %s", renderer.lastTemplate)
	}
	if renderer.lastProps["companyName"] != mail.BrandName {
		t.Fatalf("companyName %v", renderer.lastProps["companyName"])
	}
	if sender.last.Subject != "Verify your Go Boilerplate email" {
		t.Fatalf("subject %s", sender.last.Subject)
	}
}
