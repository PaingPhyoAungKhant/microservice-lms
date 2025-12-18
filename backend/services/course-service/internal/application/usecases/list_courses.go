package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type ListCoursesInput struct {
	SearchQuery   *string
	CategoryID    *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *repositories.SortDirection
}

type ListCoursesOutput struct {
	Courses []dtos.CourseDTO `json:"courses"`
	Total   int              `json:"total"`
}

type ListCoursesUseCase struct {
	courseRepo         repositories.CourseRepository
	courseCategoryRepo repositories.CourseCategoryRepository
	categoryRepo       repositories.CategoryRepository
	logger             *logger.Logger
	apiGatewayURL      string
}

func NewListCoursesUseCase(
	courseRepo repositories.CourseRepository,
	courseCategoryRepo repositories.CourseCategoryRepository,
	categoryRepo repositories.CategoryRepository,
	logger *logger.Logger,
	apiGatewayURL string,
) *ListCoursesUseCase {
	return &ListCoursesUseCase{
		courseRepo:         courseRepo,
		courseCategoryRepo: courseCategoryRepo,
		categoryRepo:       categoryRepo,
		logger:             logger,
		apiGatewayURL:      apiGatewayURL,
	}
}

func (uc *ListCoursesUseCase) Execute(ctx context.Context, input ListCoursesInput) (*ListCoursesOutput, error) {
	query := repositories.CourseQuery{
		SearchQuery:   input.SearchQuery,
		CategoryID:    input.CategoryID,
		Limit:         input.Limit,
		Offset:        input.Offset,
		SortColumn:    input.SortColumn,
		SortDirection: input.SortDirection,
	}

	result, err := uc.courseRepo.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	courseDTOs := make([]dtos.CourseDTO, len(result.Courses))
	for i, course := range result.Courses {
		courseDTOs[i].FromEntity(course, uc.apiGatewayURL)

	
		var categories []dtos.CategoryDTO
		courseCategories, err := uc.courseCategoryRepo.FindByCourseID(ctx, course.ID)
		if err != nil {
			uc.logger.Error("failed to fetch course categories", zap.String("course_id", course.ID), zap.Error(err))
		} else {
			for _, courseCategory := range courseCategories {
				category, err := uc.categoryRepo.FindByID(ctx, courseCategory.CategoryID)
				if err == nil && category != nil {
					var categoryDTO dtos.CategoryDTO
					categoryDTO.FromEntity(category)
					categories = append(categories, categoryDTO)
				}
			}
		}
		courseDTOs[i].Categories = categories
	}

	return &ListCoursesOutput{
		Courses: courseDTOs,
		Total:   result.Total,
	}, nil
}

