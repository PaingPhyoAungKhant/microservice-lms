package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var (
	ErrSectionModuleNotFound = errors.New("section module not found")
)

type UpdateSectionModuleUseCase struct {
	moduleRepo repositories.SectionModuleRepository
	publisher  messaging.Publisher
	logger     *logger.Logger
}

func NewUpdateSectionModuleUseCase(
	moduleRepo repositories.SectionModuleRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *UpdateSectionModuleUseCase {
	return &UpdateSectionModuleUseCase{
		moduleRepo: moduleRepo,
		publisher:  publisher,
		logger:     logger,
	}
}

func (uc *UpdateSectionModuleUseCase) Execute(ctx context.Context, moduleID string, input dtos.UpdateSectionModuleInput) (*dtos.SectionModuleDTO, error) {
	module, err := uc.moduleRepo.FindByID(ctx, moduleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, ErrSectionModuleNotFound
	}

	module.Update(input.Name, input.Description, input.Order)

	if err := uc.moduleRepo.Update(ctx, module); err != nil {
		return nil, err
	}

	event := events.SectionModuleUpdatedEvent{
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
		if err := uc.publisher.Publish(ctx, events.EventTypeSectionModuleUpdated, event); err != nil {
			uc.logger.Error("Failed to publish section module updated event", zap.Error(err))
		}
	}

	var dto dtos.SectionModuleDTO
	dto.FromEntity(module)
	return &dto, nil
}

