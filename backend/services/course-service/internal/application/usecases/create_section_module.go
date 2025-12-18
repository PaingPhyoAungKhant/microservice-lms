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

type CreateSectionModuleUseCase struct {
	moduleRepo  repositories.SectionModuleRepository
	sectionRepo repositories.CourseSectionRepository
	publisher   messaging.Publisher
	logger      *logger.Logger
}

func NewCreateSectionModuleUseCase(
	moduleRepo repositories.SectionModuleRepository,
	sectionRepo repositories.CourseSectionRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *CreateSectionModuleUseCase {
	return &CreateSectionModuleUseCase{
		moduleRepo:  moduleRepo,
		sectionRepo: sectionRepo,
		publisher:   publisher,
		logger:      logger,
	}
}

func (uc *CreateSectionModuleUseCase) Execute(ctx context.Context, sectionID string, input dtos.CreateSectionModuleInput) (*dtos.SectionModuleDTO, error) {
	section, err := uc.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrCourseSectionNotFound
	}

	contentType := entities.ContentType(input.ContentType)
	module := entities.NewSectionModule(sectionID, input.Name, input.Description, contentType, input.Order)

	if err := uc.moduleRepo.Create(ctx, module); err != nil {
		return nil, err
	}

	event := events.SectionModuleCreatedEvent{
		ID:              module.ID,
		CourseSectionID: module.CourseSectionID,
		ContentID:       module.ContentID,
		Name:            module.Name,
		Description:     module.Description,
		ContentType:     string(module.ContentType),
		ContentStatus:   string(module.ContentStatus),
		Order:           module.Order,
		CreatedAt:       module.CreatedAt,
		UpdatedAt:       module.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeSectionModuleCreated, event); err != nil {
			uc.logger.Error("Failed to publish section module created event", zap.Error(err))
		}
	}

	var dto dtos.SectionModuleDTO
	dto.FromEntity(module)
	return &dto, nil
}

