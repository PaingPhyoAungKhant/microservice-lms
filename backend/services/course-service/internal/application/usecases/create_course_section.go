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

type CreateCourseSectionUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	offeringRepo repositories.CourseOfferingRepository
	publisher    messaging.Publisher
	logger       *logger.Logger
}

func NewCreateCourseSectionUseCase(
	sectionRepo repositories.CourseSectionRepository,
	offeringRepo repositories.CourseOfferingRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *CreateCourseSectionUseCase {
	return &CreateCourseSectionUseCase{
		sectionRepo: sectionRepo,
		offeringRepo: offeringRepo,
		publisher:    publisher,
		logger:       logger,
	}
}

func (uc *CreateCourseSectionUseCase) Execute(ctx context.Context, offeringID string, input dtos.CreateCourseSectionInput) (*dtos.CourseSectionDTO, error) {
	offering, err := uc.offeringRepo.FindByID(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, ErrCourseOfferingNotFound
	}

	section := entities.NewCourseSection(offeringID, input.Name, input.Description, input.Order)

	if err := uc.sectionRepo.Create(ctx, section); err != nil {
		return nil, err
	}

	event := events.CourseSectionCreatedEvent{
		ID:               section.ID,
		CourseOfferingID: section.CourseOfferingID,
		Name:             section.Name,
		Description:      section.Description,
		Order:            section.Order,
		Status:           string(section.Status),
		CreatedAt:        section.CreatedAt,
		UpdatedAt:        section.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseSectionCreated, event); err != nil {
			uc.logger.Error("Failed to publish course section created event", zap.Error(err))
		}
	}

	var dto dtos.CourseSectionDTO
	dto.FromEntity(section)
	return &dto, nil
}

