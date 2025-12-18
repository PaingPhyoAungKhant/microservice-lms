package handlers

import (
	"context"
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type CourseOfferingUpdatedHandler struct {
	enrollmentRepo repositories.EnrollmentRepository
	logger         *logger.Logger
}

func NewCourseOfferingUpdatedHandler(
	enrollmentRepo repositories.EnrollmentRepository,
	logger *logger.Logger,
) *CourseOfferingUpdatedHandler {
	return &CourseOfferingUpdatedHandler{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (h *CourseOfferingUpdatedHandler) Handle(body []byte) error {
	var event events.CourseOfferingUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal course offering updated event", zap.Error(err))
		return err
	}

	
	if err := h.enrollmentRepo.UpdateCourseOfferingName(
		context.Background(),
		event.ID,
		event.Name,
	); err != nil {
		h.logger.Error("failed to update course offering name in enrollments",
			zap.String("course_offering_id", event.ID),
			zap.String("course_offering_name", event.Name),
			zap.Error(err),
		)
		return err
	}

	h.logger.Info("updated course offering name in enrollments",
		zap.String("course_offering_id", event.ID),
		zap.String("course_offering_name", event.Name),
	)

	return nil
}

