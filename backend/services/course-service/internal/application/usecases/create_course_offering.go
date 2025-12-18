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
	ErrCourseNotFound = errors.New("course not found")
)

type CreateCourseOfferingUseCase struct {
	offeringRepo repositories.CourseOfferingRepository
	courseRepo   repositories.CourseRepository
	publisher    messaging.Publisher
	logger       *logger.Logger
}

func NewCreateCourseOfferingUseCase(
	offeringRepo repositories.CourseOfferingRepository,
	courseRepo repositories.CourseRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *CreateCourseOfferingUseCase {
	return &CreateCourseOfferingUseCase{
		offeringRepo: offeringRepo,
		courseRepo:   courseRepo,
		publisher:    publisher,
		logger:       logger,
	}
}

func (uc *CreateCourseOfferingUseCase) Execute(ctx context.Context, courseID string, input dtos.CreateCourseOfferingInput) (*dtos.CourseOfferingDTO, error) {
	course, err := uc.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	offeringType := entities.OfferingType(input.OfferingType)
	offering := entities.NewCourseOffering(courseID, input.Name, input.Description, offeringType, input.Duration, input.ClassTime, input.EnrollmentCost)

	if err := uc.offeringRepo.Create(ctx, offering); err != nil {
		return nil, err
	}

	event := events.CourseOfferingCreatedEvent{
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
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseOfferingCreated, event); err != nil {
			uc.logger.Error("Failed to publish course offering created event", zap.Error(err))
		}
	}

	var dto dtos.CourseOfferingDTO
	dto.FromEntity(offering)
	return &dto, nil
}

