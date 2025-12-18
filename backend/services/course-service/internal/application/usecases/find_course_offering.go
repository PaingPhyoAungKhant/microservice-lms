package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindCourseOfferingInput struct {
	SearchQuery   *string
	CourseID      *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *repositories.SortDirection
}

type FindCourseOfferingOutput struct {
	Offerings []dtos.CourseOfferingDTO `json:"offerings"`
	Total     int                       `json:"total"`
}

type FindCourseOfferingUseCase struct {
	offeringRepo repositories.CourseOfferingRepository
	courseRepo   repositories.CourseRepository
	logger       *logger.Logger
}

func NewFindCourseOfferingUseCase(
	offeringRepo repositories.CourseOfferingRepository,
	courseRepo repositories.CourseRepository,
	logger *logger.Logger,
) *FindCourseOfferingUseCase {
	return &FindCourseOfferingUseCase{
		offeringRepo: offeringRepo,
		courseRepo:   courseRepo,
		logger:       logger,
	}
}

func (uc *FindCourseOfferingUseCase) Execute(ctx context.Context, input FindCourseOfferingInput) (*FindCourseOfferingOutput, error) {
	query := repositories.CourseOfferingQuery{
		SearchQuery:   input.SearchQuery,
		CourseID:      input.CourseID,
		Limit:         input.Limit,
		Offset:        input.Offset,
		SortColumn:    input.SortColumn,
		SortDirection: input.SortDirection,
	}

	result, err := uc.offeringRepo.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	offeringDTOs := make([]dtos.CourseOfferingDTO, len(result.Offerings))
	for i, offering := range result.Offerings {
		offeringDTOs[i].FromEntity(offering)
		
		course, err := uc.courseRepo.FindByID(ctx, offering.CourseID)
		if err == nil && course != nil {
			courseName := course.Name
			offeringDTOs[i].CourseName = &courseName
		}
	}

	return &FindCourseOfferingOutput{
		Offerings: offeringDTOs,
		Total:     result.Total,
	}, nil
}

