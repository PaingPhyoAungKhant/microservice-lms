package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type RemoveInstructorFromOfferingUseCase struct {
	instructorRepo repositories.CourseOfferingInstructorRepository
	publisher      messaging.Publisher
	logger         *logger.Logger
}

func NewRemoveInstructorFromOfferingUseCase(
	instructorRepo repositories.CourseOfferingInstructorRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *RemoveInstructorFromOfferingUseCase {
	return &RemoveInstructorFromOfferingUseCase{
		instructorRepo: instructorRepo,
		publisher:      publisher,
		logger:         logger,
	}
}

func (uc *RemoveInstructorFromOfferingUseCase) Execute(ctx context.Context, offeringID, instructorID string) error {
	instructors, err := uc.instructorRepo.FindByOfferingID(ctx, offeringID)
	if err != nil {
		return err
	}

	var instructorToRemove *entities.CourseOfferingInstructor
	for _, inst := range instructors {
		if inst.InstructorID == instructorID {
			instructorToRemove = inst
			break
		}
	}

	if instructorToRemove == nil {
		return errors.New("instructor not found in this offering")
	}

	if err := uc.instructorRepo.Delete(ctx, instructorToRemove.ID); err != nil {
		return err
	}

	event := events.InstructorRemovedFromOfferingEvent{
		ID:                 instructorToRemove.ID,
		CourseOfferingID:   instructorToRemove.CourseOfferingID,
		InstructorID:       instructorToRemove.InstructorID,
		InstructorUsername: instructorToRemove.InstructorUsername,
		RemovedAt:          time.Now().UTC(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeInstructorRemovedFromOffering, event); err != nil {
			uc.logger.Error("Failed to publish instructor removed event", zap.Error(err))
		}
	}

	return nil
}

