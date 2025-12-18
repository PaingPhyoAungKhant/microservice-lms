package handlers

import (
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/domain/templates"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/email"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type EmailVerificationRequestHandler struct {
	emailService *email.EmailService
	logger       *logger.Logger
}

func NewEmailVerificationRequestHandler(emailService *email.EmailService, logger *logger.Logger) *EmailVerificationRequestHandler {
	return &EmailVerificationRequestHandler{
		emailService: emailService,
		logger:       logger,
	}
}

func (h *EmailVerificationRequestHandler) Handle(body []byte) error {
	var event events.AuthUserRequestedEmailVerificationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal email verification request event", zap.Error(err))
		return err
	}

	templateData := map[string]interface{}{
		"VerificationURL": event.EmailVerificationURL,
	}

	htmlBody, err := h.emailService.RenderTemplate(templates.EmailVerification, templateData)
	if err != nil {
		h.logger.Error("failed to render email template", zap.Error(err))
		return err
	}

	emailData := email.EmailData{
		To:      event.Email,
		Subject: "Verify Your Email Address",
		Body:    htmlBody,
	}

	if err := h.emailService.SendEmail(emailData); err != nil {
		h.logger.Error("failed to send email verification email", zap.Error(err))
		return err
	}

	h.logger.Info("email verification request email sent",
		zap.String("user_id", event.ID),
		zap.String("email", event.Email),
	)

	return nil
}

