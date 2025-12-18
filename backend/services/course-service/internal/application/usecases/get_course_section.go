package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type GetCourseSectionInput struct {
	SectionID string
}

type GetCourseSectionUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	logger      *logger.Logger
}

func NewGetCourseSectionUseCase(
	sectionRepo repositories.CourseSectionRepository,
	logger *logger.Logger,
) *GetCourseSectionUseCase {
	return &GetCourseSectionUseCase{
		sectionRepo: sectionRepo,
		logger:      logger,
	}
}

func (uc *GetCourseSectionUseCase) Execute(ctx context.Context, input GetCourseSectionInput) (*dtos.CourseSectionDTO, error) {
	section, err := uc.sectionRepo.FindByID(ctx, input.SectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrCourseSectionNotFound
	}

	var dto dtos.CourseSectionDTO
	dto.FromEntity(section)
	return &dto, nil
}

