package handlers

import (
	"context"
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/domain/templates"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/email"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

type ForgotPasswordHandler struct {
	emailService *email.EmailService
	redis        utils.RedisInterface
	logger       *logger.Logger
}

func NewForgotPasswordHandler(emailService *email.EmailService, redis utils.RedisInterface, logger *logger.Logger) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		emailService: emailService,
		redis:        redis,
		logger:       logger,
	}
}

func (h *ForgotPasswordHandler) Handle(body []byte) error {
	var event events.AuthUserForgotPasswordEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal forgot password event", zap.Error(err))
		return err
	}

	otp := utils.GenerateOTP(6)

	ctx := context.Background()
	if err := h.redis.StoreForgotPasswordOTP(ctx, event.ID, otp); err != nil {
		h.logger.Error("failed to store forgot password OTP", zap.Error(err))
		return err
	}

	templateData := map[string]interface{}{
		"OTP": otp,
	}

	htmlBody, err := h.emailService.RenderTemplate(templates.ForgotPasswordOTP, templateData)
	if err != nil {
		h.logger.Error("failed to render email template", zap.Error(err))
		return err
	}

	emailData := email.EmailData{
		To:      event.Email,
		Subject: "Password Reset OTP",
		Body:    htmlBody,
	}

	if err := h.emailService.SendEmail(emailData); err != nil {
		h.logger.Error("failed to send forgot password OTP email", zap.Error(err))
		return err
	}

	h.logger.Info("forgot password OTP email sent",
		zap.String("user_id", event.ID),
		zap.String("email", event.Email),
	)

	return nil
}

