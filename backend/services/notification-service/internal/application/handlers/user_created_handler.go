package handlers

import (
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/domain/templates"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/email"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type UserCreatedHandler struct {
	emailService *email.EmailService
	logger       *logger.Logger
}

func NewUserCreatedHandler(emailService *email.EmailService, logger *logger.Logger) *UserCreatedHandler {
	return &UserCreatedHandler{
		emailService: emailService,
		logger:       logger,
	}
}

func (h *UserCreatedHandler) Handle(body []byte) error {
	var event events.UserCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal user created event", zap.Error(err))
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

	h.logger.Info("user created email sent",
		zap.String("user_id", event.ID),
		zap.String("email", event.Email),
	)

	return nil
}
