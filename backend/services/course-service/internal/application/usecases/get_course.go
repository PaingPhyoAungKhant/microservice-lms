package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type GetCourseInput struct {
	CourseID string
}

type GetCourseUseCase struct {
	courseRepo         repositories.CourseRepository
	courseCategoryRepo repositories.CourseCategoryRepository
	categoryRepo       repositories.CategoryRepository
	logger             *logger.Logger
	apiGatewayURL      string
}

func NewGetCourseUseCase(
	courseRepo repositories.CourseRepository,
	courseCategoryRepo repositories.CourseCategoryRepository,
	categoryRepo repositories.CategoryRepository,
	logger *logger.Logger,
	apiGatewayURL string,
) *GetCourseUseCase {
	return &GetCourseUseCase{
		courseRepo:         courseRepo,
		courseCategoryRepo: courseCategoryRepo,
		categoryRepo:       categoryRepo,
		logger:             logger,
		apiGatewayURL:      apiGatewayURL,
	}
}

func (uc *GetCourseUseCase) Execute(ctx context.Context, input GetCourseInput) (*dtos.CourseDTO, error) {
	course, err := uc.courseRepo.FindByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

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

	var dto dtos.CourseDTO
	dto.FromEntity(course, uc.apiGatewayURL)
	dto.Categories = categories
	return &dto, nil
}

