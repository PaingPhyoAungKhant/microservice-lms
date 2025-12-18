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

type UpdateCourseInput struct {
	CourseID string
}

type UpdateCourseUseCase struct {
	courseRepo         repositories.CourseRepository
	courseCategoryRepo repositories.CourseCategoryRepository
	categoryRepo       repositories.CategoryRepository
	publisher          messaging.Publisher
	logger             *logger.Logger
	apiGatewayURL      string
}

func NewUpdateCourseUseCase(
	courseRepo repositories.CourseRepository,
	courseCategoryRepo repositories.CourseCategoryRepository,
	categoryRepo repositories.CategoryRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
	apiGatewayURL string,
) *UpdateCourseUseCase {
	return &UpdateCourseUseCase{
		courseRepo:         courseRepo,
		courseCategoryRepo: courseCategoryRepo,
		categoryRepo:       categoryRepo,
		publisher:          publisher,
		logger:             logger,
		apiGatewayURL:      apiGatewayURL,
	}
}

func (uc *UpdateCourseUseCase) Execute(ctx context.Context, input UpdateCourseInput, updateData dtos.UpdateCourseInput) (*dtos.CourseDTO, error) {
	course, err := uc.courseRepo.FindByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	course.Update(updateData.Name, updateData.Description, updateData.ThumbnailID)
	if err := uc.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	if updateData.CategoryIDs != nil {
		currentCourseCategories, err := uc.courseCategoryRepo.FindByCourseID(ctx, course.ID)
		if err != nil {
			uc.logger.Error("failed to fetch current course categories", zap.String("course_id", course.ID), zap.Error(err))
			return nil, fmt.Errorf("failed to fetch current categories: %w", err)
		}

		currentCategoryIDs := make(map[string]bool)
		for _, cc := range currentCourseCategories {
			currentCategoryIDs[cc.CategoryID] = true
		}

		newCategoryIDs := make(map[string]bool)
		for _, categoryID := range updateData.CategoryIDs {
			newCategoryIDs[categoryID] = true
		}

		for _, cc := range currentCourseCategories {
			if !newCategoryIDs[cc.CategoryID] {
				if err := uc.courseCategoryRepo.Delete(ctx, course.ID, cc.CategoryID); err != nil {
					uc.logger.Error("failed to delete course category", zap.String("course_id", course.ID), zap.String("category_id", cc.CategoryID), zap.Error(err))
					return nil, fmt.Errorf("failed to remove category association: %w", err)
				}
			}
		}

		for _, categoryID := range updateData.CategoryIDs {
			if !currentCategoryIDs[categoryID] {
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

	event := events.CourseUpdatedEvent{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		ThumbnailID: course.ThumbnailID,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseUpdated, event); err != nil {
			uc.logger.Error("Failed to publish course updated event", zap.Error(err))
		}
	}

	var dto dtos.CourseDTO
	dto.FromEntity(course, uc.apiGatewayURL)
	dto.Categories = categories
	return &dto, nil
}

