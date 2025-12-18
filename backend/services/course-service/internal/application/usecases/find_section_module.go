package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindSectionModuleInput struct {
	SectionID string
}

type FindSectionModuleOutput struct {
	Modules []dtos.SectionModuleDTO `json:"modules"`
	Total   int                      `json:"total"`
}

type FindSectionModuleUseCase struct {
	moduleRepo repositories.SectionModuleRepository
	logger     *logger.Logger
}

func NewFindSectionModuleUseCase(
	moduleRepo repositories.SectionModuleRepository,
	logger *logger.Logger,
) *FindSectionModuleUseCase {
	return &FindSectionModuleUseCase{
		moduleRepo: moduleRepo,
		logger:     logger,
	}
}

func (uc *FindSectionModuleUseCase) Execute(ctx context.Context, input FindSectionModuleInput) (*FindSectionModuleOutput, error) {
	modules, err := uc.moduleRepo.FindBySectionID(ctx, input.SectionID)
	if err != nil {
		return nil, err
	}

	moduleDTOs := make([]dtos.SectionModuleDTO, len(modules))
	for i, module := range modules {
		moduleDTOs[i].FromEntity(module)
	}

	return &FindSectionModuleOutput{
		Modules: moduleDTOs,
		Total:   len(moduleDTOs),
	}, nil
}

