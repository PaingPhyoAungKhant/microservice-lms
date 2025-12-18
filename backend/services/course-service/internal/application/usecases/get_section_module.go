package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type GetSectionModuleInput struct {
	ModuleID string
}

type GetSectionModuleUseCase struct {
	moduleRepo repositories.SectionModuleRepository
	logger     *logger.Logger
}

func NewGetSectionModuleUseCase(
	moduleRepo repositories.SectionModuleRepository,
	logger *logger.Logger,
) *GetSectionModuleUseCase {
	return &GetSectionModuleUseCase{
		moduleRepo: moduleRepo,
		logger:     logger,
	}
}

func (uc *GetSectionModuleUseCase) Execute(ctx context.Context, input GetSectionModuleInput) (*dtos.SectionModuleDTO, error) {
	module, err := uc.moduleRepo.FindByID(ctx, input.ModuleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, ErrSectionModuleNotFound
	}

	var dto dtos.SectionModuleDTO
	dto.FromEntity(module)
	return &dto, nil
}

