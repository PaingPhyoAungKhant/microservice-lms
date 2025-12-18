package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var (
	ErrCourseOfferingNotFound = errors.New("course offering not found")
)

type UpdateCourseOfferingUseCase struct {
	offeringRepo repositories.CourseOfferingRepository
	publisher    messaging.Publisher
	logger       *logger.Logger
}

func NewUpdateCourseOfferingUseCase(
	offeringRepo repositories.CourseOfferingRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *UpdateCourseOfferingUseCase {
	return &UpdateCourseOfferingUseCase{
		offeringRepo: offeringRepo,
		publisher:    publisher,
		logger:       logger,
	}
}

func (uc *UpdateCourseOfferingUseCase) Execute(ctx context.Context, offeringID string, input dtos.UpdateCourseOfferingInput) (*dtos.CourseOfferingDTO, error) {
	offering, err := uc.offeringRepo.FindByID(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, ErrCourseOfferingNotFound
	}

	offeringType := entities.OfferingType(input.OfferingType)
	offering.Update(input.Name, input.Description, offeringType, input.Duration, input.ClassTime, input.EnrollmentCost)
	
	if input.Status != nil {
		status := entities.OfferingStatus(*input.Status)
		offering.UpdateStatus(status)
	}

	if err := uc.offeringRepo.Update(ctx, offering); err != nil {
		return nil, err
	}

	event := events.CourseOfferingUpdatedEvent{
		ID:             offering.ID,
		CourseID:       offering.CourseID,
		Name:           offering.Name,
		Description:    offering.Description,
		OfferingType:   string(offering.OfferingType),
		Status:         string(offering.Status),
		Duration:       offering.Duration,
		ClassTime:      offering.ClassTime,
		EnrollmentCost: offering.EnrollmentCost,
		CreatedAt:      offering.CreatedAt,
		UpdatedAt:      offering.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseOfferingUpdated, event); err != nil {
			uc.logger.Error("Failed to publish course offering updated event", zap.Error(err))
		}
	}

	var dto dtos.CourseOfferingDTO
	dto.FromEntity(offering)
	return &dto, nil
}

