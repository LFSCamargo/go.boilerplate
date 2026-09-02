package smtp_test

import (
	"context"
	"testing"

	"go.boilerplate/src/mail/smtp"
)

func TestNewClient_RequiresHostAndFrom(t *testing.T) {
	_, err := smtp.NewClient(smtp.Config{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}

	_, err = smtp.NewClient(smtp.Config{Host: "localhost", Port: "1025"})
	if err == nil {
		t.Fatal("expected error without from address")
	}
}

func TestClient_Send_ValidatesRecipient(t *testing.T) {
	client, err := smtp.NewClient(smtp.Config{
		Host:     "localhost",
		Port:     "1025",
		Security: "none",
		From:     "noreply@go.boilerplate",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Send(context.Background(), smtp.Message{HTML: "<p>hi</p>"})
	if err == nil {
		t.Fatal("expected error for missing recipient")
	}
}
