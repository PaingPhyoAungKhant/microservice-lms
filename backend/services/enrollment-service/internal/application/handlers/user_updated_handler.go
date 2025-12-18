package handlers

import (
	"context"
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type UserUpdatedHandler struct {
	enrollmentRepo repositories.EnrollmentRepository
	logger         *logger.Logger
}

func NewUserUpdatedHandler(
	enrollmentRepo repositories.EnrollmentRepository,
	logger *logger.Logger,
) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (h *UserUpdatedHandler) Handle(body []byte) error {
	var event events.UserUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal user updated event", zap.Error(err))
		return err
	}

	if err := h.enrollmentRepo.UpdateStudentUsername(
		context.Background(),
		event.ID,
		event.Username,
	); err != nil {
		h.logger.Error("failed to update student username in enrollments",
			zap.String("student_id", event.ID),
			zap.String("username", event.Username),
			zap.Error(err),
		)
		return err
	}

	h.logger.Info("updated student username in enrollments",
		zap.String("student_id", event.ID),
		zap.String("username", event.Username),
	)

	return nil
}

