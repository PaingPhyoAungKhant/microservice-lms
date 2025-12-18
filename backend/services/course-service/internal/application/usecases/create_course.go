package usecases

import (
	"context"
	"fmt"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type CreateCourseUseCase struct {
	courseRepo         repositories.CourseRepository
	courseCategoryRepo repositories.CourseCategoryRepository
	categoryRepo       repositories.CategoryRepository
	publisher          messaging.Publisher
	logger             *logger.Logger
	apiGatewayURL      string
}

func NewCreateCourseUseCase(
	courseRepo repositories.CourseRepository,
	courseCategoryRepo repositories.CourseCategoryRepository,
	categoryRepo repositories.CategoryRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
	apiGatewayURL string,
) *CreateCourseUseCase {
	return &CreateCourseUseCase{
		courseRepo:         courseRepo,
		courseCategoryRepo: courseCategoryRepo,
		categoryRepo:       categoryRepo,
		publisher:          publisher,
		logger:             logger,
		apiGatewayURL:      apiGatewayURL,
	}
}

func (uc *CreateCourseUseCase) Execute(ctx context.Context, input dtos.CreateCourseInput) (*dtos.CourseDTO, error) {
	course := entities.NewCourse(input.Name, input.Description, input.ThumbnailID)
	if err := uc.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	if len(input.CategoryIDs) > 0 {
		for _, categoryID := range input.CategoryIDs {
			category, err := uc.categoryRepo.FindByID(ctx, categoryID)
			if err != nil {
				uc.logger.Error("failed to find category", zap.String("category_id", categoryID), zap.Error(err))
				return nil, fmt.Errorf("category not found: %s", categoryID)
			}
			if category == nil {
				return nil, fmt.Errorf("category not found: %s", categoryID)
			}

			courseCategory := entities.NewCourseCategory(course.ID, categoryID)
			if err := uc.courseCategoryRepo.Create(ctx, courseCategory); err != nil {
				uc.logger.Error("failed to create course category", zap.String("course_id", course.ID), zap.String("category_id", categoryID), zap.Error(err))
				return nil, fmt.Errorf("failed to associate category with course: %w", err)
			}
		}
	}

	var categories []dtos.CategoryDTO
	if len(input.CategoryIDs) > 0 {
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
	}

	event := events.CourseCreatedEvent{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		ThumbnailID: course.ThumbnailID,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseCreated, event); err != nil {
			uc.logger.Error("Failed to publish course created event", zap.Error(err))
		}
	}

	var dto dtos.CourseDTO
	dto.FromEntity(course, uc.apiGatewayURL)
	dto.Categories = categories
	return &dto, nil
}

