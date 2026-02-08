package adapter

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/aiagent/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
)

// EmailProvider defines the interface for sending emails
type EmailProvider interface {
	Send(ctx context.Context, to []string, subject string, htmlBody string, textBody string) error
}

// smtpAdapter implements EmailProvider using net/smtp
type smtpAdapter struct {
	cfg config.EmailConfig
}

// NewSMTPAdapter creates a new SMTP adapter instance
func NewSMTPAdapter(cfg config.EmailConfig) EmailProvider {
	return &smtpAdapter{
		cfg: cfg,
	}
}

// Send sends an email via SMTP
func (a *smtpAdapter) Send(ctx context.Context, to []string, subject string, htmlBody string, textBody string) error {
	if len(to) == 0 {
		return nil
	}

	// Prepare authentication
	auth := smtp.PlainAuth("", a.cfg.User, a.cfg.Password, a.cfg.Host)

	// Prepare message
	delimiter := "**NextPart**"
	header := make(map[string]string)
	header["From"] = a.cfg.From
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("multipart/alternative; boundary=\"%s\"", delimiter)

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	// Text body
	msg.WriteString(fmt.Sprintf("--%s\r\n", delimiter))
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(textBody)
	msg.WriteString("\r\n")

	// HTML body
	msg.WriteString(fmt.Sprintf("--%s\r\n", delimiter))
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s--", delimiter))

	addr := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)

	// Send the email
	// Note: smtp.SendMail doesn't take a context, so we just run it.
	// In a more robust implementation, we might want to use a dialer that supports context.
	err := smtp.SendMail(addr, auth, a.cfg.From, to, []byte(msg.String()))
	if err != nil {
		log.Error().Err(err).
			Str("subject", subject).
			Strs("to", to).
			Msg("failed to send email via SMTP")
		return err
	}

	log.Info().
		Str("subject", subject).
		Strs("to", to).
		Msg("email sent successfully")

	return nil
}
