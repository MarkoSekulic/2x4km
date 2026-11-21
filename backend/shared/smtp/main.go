package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

// Config holds the SMTP server connection details.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Client is an SMTP mailer instance.
type Client struct {
	config Config
	auth   smtp.Auth
}

// Message represents an email to be sent.
type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// New creates a new SMTP mailer Client.
func New(cfg Config) *Client {
	var auth smtp.Auth
	if cfg.Username != "" && cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return &Client{
		config: cfg,
		auth:   auth,
	}
}

// Send formats and dispatches a plain text or raw HTML email.
func (c *Client) Send(msg Message) error {
	if msg.From == "" {
		msg.From = c.config.From
	}
	if msg.From == "" {
		return fmt.Errorf("mailer: sender address (From) is required")
	}
	if len(msg.To) == 0 {
		return fmt.Errorf("mailer: at least one recipient is required")
	}
	if msg.Subject == "" {
		return fmt.Errorf("mailer: subject is required")
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", msg.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	if msg.IsHTML {
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		buf.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}

	buf.WriteString("\r\n")
	buf.WriteString(msg.Body)

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	if err := smtp.SendMail(addr, nil, msg.From, msg.To, buf.Bytes()); err != nil {
		return fmt.Errorf("mailer: failed to send email: %w", err)
	}

	return nil
}

// SendTemplate parses an HTML template string, injects data, and sends the email.
func (c *Client) SendTemplate(to []string, subject, tmplStr string, data any) error {
	tmpl, err := template.New("email").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("mailer: failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("mailer: failed to execute template: %w", err)
	}

	return c.Send(Message{
		To:      to,
		Subject: subject,
		Body:    body.String(),
		IsHTML:  true,
	})
}

// SendTemplateFile loads an HTML file from disk, injects data, and sends the email.
func (c *Client) SendTemplateFile(to []string, subject, filePath string, data any) error {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return fmt.Errorf("mailer: failed to load template file: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("mailer: failed to execute template file: %w", err)
	}

	return c.Send(Message{
		To:      to,
		Subject: subject,
		Body:    body.String(),
		IsHTML:  true,
	})
}
