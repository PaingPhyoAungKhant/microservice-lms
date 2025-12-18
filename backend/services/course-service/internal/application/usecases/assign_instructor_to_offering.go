package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type AssignInstructorToOfferingUseCase struct {
	instructorRepo repositories.CourseOfferingInstructorRepository
	offeringRepo   repositories.CourseOfferingRepository
	publisher      messaging.Publisher
	logger         *logger.Logger
}

func NewAssignInstructorToOfferingUseCase(
	instructorRepo repositories.CourseOfferingInstructorRepository,
	offeringRepo repositories.CourseOfferingRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *AssignInstructorToOfferingUseCase {
	return &AssignInstructorToOfferingUseCase{
		instructorRepo: instructorRepo,
		offeringRepo:   offeringRepo,
		publisher:      publisher,
		logger:         logger,
	}
}

func (uc *AssignInstructorToOfferingUseCase) Execute(ctx context.Context, offeringID string, input dtos.AssignInstructorInput, instructorUsername string) (*dtos.CourseOfferingInstructorDTO, error) {	
	offering, err := uc.offeringRepo.FindByID(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, ErrCourseOfferingNotFound
	}

	instructor := entities.NewCourseOfferingInstructor(offeringID, input.InstructorID, instructorUsername)

	if err := uc.instructorRepo.Create(ctx, instructor); err != nil {
		return nil, err
	}

	event := events.InstructorAssignedToOfferingEvent{
		ID:                 instructor.ID,
		CourseOfferingID:   instructor.CourseOfferingID,
		InstructorID:       instructor.InstructorID,
		InstructorUsername: instructor.InstructorUsername,
		CreatedAt:         instructor.CreatedAt,
		UpdatedAt:         instructor.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeInstructorAssignedToOffering, event); err != nil {
			uc.logger.Error("Failed to publish instructor assigned event", zap.Error(err))
		}
	}

	var dto dtos.CourseOfferingInstructorDTO
	dto.FromEntity(instructor)
	return &dto, nil
}

