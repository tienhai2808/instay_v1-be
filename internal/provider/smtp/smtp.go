package smtp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/types"
)

//go:embed templates/auth.html
var authTemplate embed.FS

type SMTPProvider interface {
	Send(to, subject, body string) error
	
	AuthEmail(to, subject, otp string) error
}

type smtpProviderImpl struct {
	cfg  *config.Config
	auth smtp.Auth
}

func NewSMTPProvider(cfg *config.Config) SMTPProvider {
	auth := smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Password, cfg.SMTP.Host)
	return &smtpProviderImpl{
		cfg,
		auth,
	}
}

func (s *smtpProviderImpl) Send(to, subject, body string) error {
	msg := fmt.Appendf(nil, "Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", subject, body)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTP.Host, s.cfg.SMTP.Port)
	return smtp.SendMail(addr, s.auth, s.cfg.SMTP.User, []string{to}, msg)
}

func (s *smtpProviderImpl) AuthEmail(to, subject, otp string) error {
	tmpl, err := template.ParseFS(authTemplate, "templates/auth.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := types.AuthEmailData{
		Subject: subject,
		Otp:     otp,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return s.Send(to, subject, body.String())
}
