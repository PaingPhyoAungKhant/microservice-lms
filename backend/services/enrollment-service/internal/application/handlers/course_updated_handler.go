package handlers

import (
	"context"
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type CourseUpdatedHandler struct {
	enrollmentRepo repositories.EnrollmentRepository
	logger         *logger.Logger
}

func NewCourseUpdatedHandler(
	enrollmentRepo repositories.EnrollmentRepository,
	logger *logger.Logger,
) *CourseUpdatedHandler {
	return &CourseUpdatedHandler{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (h *CourseUpdatedHandler) Handle(body []byte) error {
	var event events.CourseUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal course updated event", zap.Error(err))
		return err
	}

	// Update course name in all enrollments for this course
	if err := h.enrollmentRepo.UpdateCourseName(
		context.Background(),
		event.ID,
		event.Name,
	); err != nil {
		h.logger.Error("failed to update course name in enrollments",
			zap.String("course_id", event.ID),
			zap.String("course_name", event.Name),
			zap.Error(err),
		)
		return err
	}

	h.logger.Info("updated course name in enrollments",
		zap.String("course_id", event.ID),
		zap.String("course_name", event.Name),
	)

	return nil
}

