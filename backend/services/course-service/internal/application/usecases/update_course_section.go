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
	ErrCourseSectionNotFound = errors.New("course section not found")
)

type UpdateCourseSectionUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	publisher   messaging.Publisher
	logger      *logger.Logger
}

func NewUpdateCourseSectionUseCase(
	sectionRepo repositories.CourseSectionRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *UpdateCourseSectionUseCase {
	return &UpdateCourseSectionUseCase{
		sectionRepo: sectionRepo,
		publisher:   publisher,
		logger:      logger,
	}
}

func (uc *UpdateCourseSectionUseCase) Execute(ctx context.Context, sectionID string, input dtos.UpdateCourseSectionInput) (*dtos.CourseSectionDTO, error) {
	section, err := uc.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrCourseSectionNotFound
	}

	section.Update(input.Name, input.Description, input.Order)

	if input.Status != nil {
		status := entities.SectionStatus(*input.Status)
		section.UpdateStatus(status)
	}

	if err := uc.sectionRepo.Update(ctx, section); err != nil {
		return nil, err
	}

	event := events.CourseSectionUpdatedEvent{
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
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseSectionUpdated, event); err != nil {
			uc.logger.Error("Failed to publish course section updated event", zap.Error(err))
		}
	}

	var dto dtos.CourseSectionDTO
	dto.FromEntity(section)
	return &dto, nil
}

	