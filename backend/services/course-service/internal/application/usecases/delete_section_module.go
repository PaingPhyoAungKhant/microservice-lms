package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type DeleteSectionModuleInput struct {
	ModuleID string
}

type DeleteSectionModuleOutput struct {
	Message string
}

type DeleteSectionModuleUseCase struct {
	moduleRepo repositories.SectionModuleRepository
	publisher  messaging.Publisher
	logger     *logger.Logger
}

func NewDeleteSectionModuleUseCase(
	moduleRepo repositories.SectionModuleRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteSectionModuleUseCase {
	return &DeleteSectionModuleUseCase{
		moduleRepo: moduleRepo,
		publisher:  publisher,
		logger:     logger,
	}
}

func (uc *DeleteSectionModuleUseCase) Execute(ctx context.Context, input DeleteSectionModuleInput) (*DeleteSectionModuleOutput, error) {
	module, err := uc.moduleRepo.FindByID(ctx, input.ModuleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, ErrSectionModuleNotFound
	}

	if err := uc.moduleRepo.Delete(ctx, input.ModuleID); err != nil {
		return nil, fmt.Errorf("failed to delete section module: %w", err)
	}

	event := events.SectionModuleDeletedEvent{
		ID:        module.ID,
		DeletedAt: time.Now(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeSectionModuleDeleted, event); err != nil {
			uc.logger.Error("failed to publish section module deleted event", zap.Error(err))
		}
	}

	uc.logger.Info("Section module deleted successfully",
		zap.String("module_id", module.ID),
	)

	return &DeleteSectionModuleOutput{
		Message: "Section module deleted successfully",
	}, nil
}

