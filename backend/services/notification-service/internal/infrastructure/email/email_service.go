package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/config"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type EmailService struct {
	config *config.SMTPConfig
	logger *logger.Logger
}

func NewEmailService(config *config.SMTPConfig, logger *logger.Logger) *EmailService {
	return &EmailService{
		config: config,
		logger: logger,
	}
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func (s *EmailService) SendEmail(data EmailData) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	to := []string{data.To}


	s.logger.Info("sending email",
		zap.String("to", data.To),
		zap.String("subject", data.Subject),
		zap.Int("body_length", len(data.Body)),
	)

	msg := []byte(fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", data.To) +
		fmt.Sprintf("Subject: %s\r\n", data.Subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		data.Body + "\r\n")

	err := smtp.SendMail(addr, auth, s.config.FromEmail, to, msg)
	if err != nil {
		s.logger.Error("failed to send email",
			zap.String("to", data.To),
			zap.String("subject", data.Subject),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("email sent successfully",
		zap.String("to", data.To),
		zap.String("subject", data.Subject),
	)

	return nil
}

func (s *EmailService) RenderTemplate(templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

