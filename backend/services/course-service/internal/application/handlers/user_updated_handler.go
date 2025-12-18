package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type UserUpdatedHandler struct {
	instructorRepo repositories.CourseOfferingInstructorRepository
	logger         *logger.Logger
}

func NewUserUpdatedHandler(
	instructorRepo repositories.CourseOfferingInstructorRepository,
	logger *logger.Logger,
) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		instructorRepo: instructorRepo,
		logger:         logger,
	}
}

func (h *UserUpdatedHandler) Handle(body []byte) error {
	ctx := context.Background()
	
	var event events.UserUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal user updated event", zap.Error(err))
		return err
	}

	instructors, err := h.instructorRepo.FindByInstructorID(ctx, event.ID)
	if err != nil {
		h.logger.Error("failed to find instructors", zap.Error(err), zap.String("instructor_id", event.ID))
		return err
	}
	
	for _, instructor := range instructors {
		if instructor.InstructorUsername != event.Username {
			instructor.UpdateInstructorUsername(event.Username)
			instructor.UpdatedAt = time.Now().UTC()
			if err := h.instructorRepo.Update(ctx, instructor); err != nil {
				h.logger.Error("failed to update instructor username",
					zap.Error(err),
					zap.String("instructor_id", instructor.ID),
					zap.String("new_username", event.Username),
				)
				return err
			}
			h.logger.Info("updated instructor username",
				zap.String("instructor_id", instructor.ID),
				zap.String("old_username", instructor.InstructorUsername),
				zap.String("new_username", event.Username),
			)
		}
	}

	return nil
}

