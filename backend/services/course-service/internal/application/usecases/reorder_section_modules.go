package usecases

import (
	"context"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type ReorderSectionModulesInput struct {
	SectionID string        `json:"section_id"`
	Items     []ReorderItem `json:"items"`
}

type ReorderSectionModulesUseCase struct {
	moduleRepo repositories.SectionModuleRepository
	logger      *logger.Logger
}

func NewReorderSectionModulesUseCase(
	moduleRepo repositories.SectionModuleRepository,
	logger *logger.Logger,
) *ReorderSectionModulesUseCase {
	return &ReorderSectionModulesUseCase{
		moduleRepo: moduleRepo,
		logger:     logger,
	}
}

func (uc *ReorderSectionModulesUseCase) Execute(ctx context.Context, input ReorderSectionModulesInput) error {
	if input.SectionID == "" {
		return ErrInvalidReorderInput
	}
	if len(input.Items) == 0 {
		return ErrInvalidReorderInput
	}

	for _, item := range input.Items {
		module, err := uc.moduleRepo.FindByID(ctx, item.ID)
		if err != nil {
			return err
		}
		if module == nil {
			return ErrSectionModuleNotFound
		}
		if module.CourseSectionID != input.SectionID {
			return ErrInvalidReorderInput
		}
	}

	for _, item := range input.Items {
		module, err := uc.moduleRepo.FindByID(ctx, item.ID)
		if err != nil {
			return err
		}
		if module == nil {
			continue
		}

		module.Order = item.Order
		module.UpdatedAt = time.Now().UTC()

		if err := uc.moduleRepo.Update(ctx, module); err != nil {
			uc.logger.Error("Failed to update module order", zap.String("module_id", item.ID), zap.Error(err))
			return err
		}
	}

	uc.logger.Info("Section modules reordered successfully", zap.String("section_id", input.SectionID))
	return nil
}

