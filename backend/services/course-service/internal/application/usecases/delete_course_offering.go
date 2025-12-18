package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type DeleteCourseOfferingUseCase struct {
	offeringRepo repositories.CourseOfferingRepository
	publisher    messaging.Publisher
	logger       *logger.Logger
}

func NewDeleteCourseOfferingUseCase(
	offeringRepo repositories.CourseOfferingRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteCourseOfferingUseCase {
	return &DeleteCourseOfferingUseCase{
		offeringRepo: offeringRepo,
		publisher:    publisher,
		logger:       logger,
	}
}

func (uc *DeleteCourseOfferingUseCase) Execute(ctx context.Context, offeringID string) error {
	offering, err := uc.offeringRepo.FindByID(ctx, offeringID)
	if err != nil {
		return err
	}
	if offering == nil {
		return ErrCourseOfferingNotFound
	}

	if err := uc.offeringRepo.Delete(ctx, offeringID); err != nil {
		return err
	}

	event := events.CourseOfferingDeletedEvent{
		ID:        offering.ID,
		DeletedAt: offering.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseOfferingDeleted, event); err != nil {
			uc.logger.Error("Failed to publish course offering deleted event", zap.Error(err))
		}
	}

	return nil
}

