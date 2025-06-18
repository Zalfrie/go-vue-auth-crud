package services

import (
	"fmt"
	"net/smtp"

	"go-vue-auth-crud/config"
)

// EmailService handles email sending
type EmailService struct {
	smtpHost string
	smtpPort int
	user     string
	pass     string
}

// NewEmailService constructor
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		smtpHost: cfg.SMTPHost,
		smtpPort: cfg.SMTPPort,
		user:     cfg.SMTPUser,
		pass:     cfg.SMTPPass,
	}
}

// SendRegistration notification
func (e *EmailService) SendRegistration(to, name string) error {
	subject := "Welcome to Go-Vue Auth App"
	body := fmt.Sprintf("Hi %s,

Thank you for registering!", name)
	return e.send(to, subject, body)
}

// SendPasswordReset email
func (e *EmailService) SendPasswordReset(to, token string) error {
	subject := "Password Reset Request"
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	body := fmt.Sprintf("Click here to reset: %s", resetLink)
	return e.send(to, subject, body)
}

// Generic send
func (e *EmailService) send(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s
To: %s
Subject: %s

%s", e.user, to, subject, body)
	auth := smtp.PlainAuth("", e.user, e.pass, e.smtpHost)
	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)
	return smtp.SendMail(addr, auth, e.user, []string{to}, []byte(msg))
}