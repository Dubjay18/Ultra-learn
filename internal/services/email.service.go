package services

import (
	"Ultra-learn/internal/logger"
	"bytes"
	"html/template"
	"net/smtp"
	"os"
)

var (
	FromEmail     = os.Getenv("FROM_EMAIL")
	EmailPassword = os.Getenv("FROM_EMAIL_PASSWORD")
	EmailHost     = os.Getenv("FROM_EMAIL_SMTP")
)

type EmailService interface {
	SendEmail(to string, subject string, body string) error
	SendSignUpEmail(to string, name string) error
}

type DefaultEmailService struct {
}

func (s *DefaultEmailService) SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", FromEmail, EmailPassword, EmailHost)
	err := smtp.SendMail(EmailHost, auth, FromEmail, []string{to}, []byte("Subject: "+subject+"\r\n"+body))
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultEmailService) SendSignUpEmail(to string, name string) error {
	// Load the email template
	t, err := template.ParseFiles("email_templates/signup_welcome.html")
	if err != nil {
		logger.Error("Error loading email template: " + err.Error())
		return err
	}
	data := struct {
		FirstName string
	}{
		FirstName: name,
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		logger.Error("Error executing email template: " + err.Error())
		return err

	}
	emailBody := tpl.String()
	if err := s.SendEmail(to, "Welcome to Our Service", emailBody); err != nil {
		logger.Error("Error sending email: " + err.Error())
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func NewEmailService() EmailService {
	return &DefaultEmailService{}
}
