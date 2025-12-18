package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/domain/templates"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/email"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type StudentRegisteredHandler struct {
	emailService *email.EmailService
	logger       *logger.Logger
}

func NewStudentRegisteredHandler(emailService *email.EmailService, logger *logger.Logger) *StudentRegisteredHandler {
	return &StudentRegisteredHandler{
		emailService: emailService,
		logger:       logger,
	}
}

func (h *StudentRegisteredHandler) Handle(body []byte) error {
	var event events.AuthStudentRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal student registered event", zap.Error(err))
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

	
	h.logger.Info("email template rendered",
		zap.String("template_length", fmt.Sprintf("%d", len(templates.EmailVerification))),
		zap.String("body_length", fmt.Sprintf("%d", len(htmlBody))),
		zap.String("verification_url", event.EmailVerificationURL),
	)

	emailData := email.EmailData{
		To:      event.Email,
		Subject: "Verify Your Email Address",
		Body:    htmlBody,
	}

	if err := h.emailService.SendEmail(emailData); err != nil {
		h.logger.Error("failed to send email verification email", zap.Error(err))
		return err
	}

	h.logger.Info("student registered email sent",
		zap.String("user_id", event.ID),
		zap.String("email", event.Email),
	)

	return nil
}

