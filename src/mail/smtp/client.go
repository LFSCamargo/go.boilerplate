package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type Message struct {
	To      string
	Subject string
	HTML    string
	Text    string
}

type Config struct {
	Host     string
	Port     string
	Security string
	User     string
	Password string
	From     string
}

type Client struct {
	host     string
	port     string
	security string
	user     string
	password string
	from     string
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("smtp host is required")
	}
	if cfg.Port == "" {
		cfg.Port = "587"
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("mail from address is required")
	}

	return &Client{
		host:     cfg.Host,
		port:     cfg.Port,
		security: strings.ToLower(strings.TrimSpace(cfg.Security)),
		user:     cfg.User,
		password: cfg.Password,
		from:     cfg.From,
	}, nil
}

func (c *Client) Send(ctx context.Context, msg Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if msg.To == "" {
		return fmt.Errorf("recipient is required")
	}
	if msg.HTML == "" && msg.Text == "" {
		return fmt.Errorf("message body is required")
	}

	body := msg.HTML
	contentType := "text/html; charset=UTF-8"
	if body == "" {
		body = msg.Text
		contentType = "text/plain; charset=UTF-8"
	}

	subject := msg.Subject
	if subject == "" {
		subject = "go.boilerplate notification"
	}

	payload := buildMessage(c.from, msg.To, subject, contentType, body)
	addr := net.JoinHostPort(c.host, c.port)

	var auth smtp.Auth
	if c.user != "" {
		auth = smtp.PlainAuth("", c.user, c.password, c.host)
	}

	switch c.security {
	case "", "none", "plain":
		return smtp.SendMail(addr, auth, c.from, []string{msg.To}, payload)
	case "starttls":
		return c.sendStartTLS(addr, auth, msg.To, payload)
	case "tls", "ssl":
		return c.sendTLS(addr, auth, msg.To, payload)
	default:
		return fmt.Errorf("unsupported smtp security mode: %q", c.security)
	}
}

func (c *Client) sendStartTLS(addr string, auth smtp.Auth, to string, payload []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: c.host}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(c.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data close: %w", err)
	}
	return client.Quit()
}

func (c *Client) sendTLS(addr string, auth smtp.Auth, to string, payload []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(c.from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func buildMessage(from, to, subject, contentType, body string) []byte {
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType,
		"Content-Transfer-Encoding: 8bit",
		"",
		body,
	}
	return []byte(strings.Join(headers, "\r\n"))
}
