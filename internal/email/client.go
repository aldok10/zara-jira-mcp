package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/aldok10/zara-jira-mcp/config"
)

// Client sends emails via SMTP.
type Client struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewClient(cfg *config.Config) *Client {
	port := cfg.Email.SMTPPort
	if port == "" {
		port = "587"
	}
	return &Client{
		host:     cfg.Email.SMTPHost,
		port:     port,
		username: cfg.Email.Username,
		password: cfg.Email.Password,
		from:     cfg.Email.From,
	}
}

func (c *Client) Available() bool {
	return c.host != "" && c.from != ""
}

// Send sends a plain text email.
func (c *Client) Send(ctx context.Context, to, subject, body string) error {
	if !c.Available() {
		return fmt.Errorf("email not configured: set EMAIL_SMTP_HOST and EMAIL_FROM")
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		c.from, to, subject, body)

	auth := smtp.PlainAuth("", c.username, c.password, c.host)
	addr := c.host + ":" + c.port
	return smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg))
}

// SendHTML sends an HTML email.
func (c *Client) SendHTML(ctx context.Context, to, subject, html string) error {
	if !c.Available() {
		return fmt.Errorf("email not configured")
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		c.from, to, subject, html)

	auth := smtp.PlainAuth("", c.username, c.password, c.host)
	addr := c.host + ":" + c.port
	return smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg))
}
