package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindCourseSectionInput struct {
	OfferingID string
}

type FindCourseSectionOutput struct {
	Sections []dtos.CourseSectionDTO `json:"sections"`
	Total    int                      `json:"total"`
}

type FindCourseSectionUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	logger      *logger.Logger
}

func NewFindCourseSectionUseCase(
	sectionRepo repositories.CourseSectionRepository,
	logger *logger.Logger,
) *FindCourseSectionUseCase {
	return &FindCourseSectionUseCase{
		sectionRepo: sectionRepo,
		logger:      logger,
	}
}

func (uc *FindCourseSectionUseCase) Execute(ctx context.Context, input FindCourseSectionInput) (*FindCourseSectionOutput, error) {
	sections, err := uc.sectionRepo.FindByOfferingID(ctx, input.OfferingID)
	if err != nil {
		return nil, err
	}

	sectionDTOs := make([]dtos.CourseSectionDTO, len(sections))
	for i, section := range sections {
		sectionDTOs[i].FromEntity(section)
	}

	return &FindCourseSectionOutput{
		Sections: sectionDTOs,
		Total:    len(sectionDTOs),
	}, nil
}

