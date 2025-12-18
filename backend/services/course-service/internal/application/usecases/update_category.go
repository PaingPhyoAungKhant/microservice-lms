package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type UpdateCategoryInput struct {
	CategoryID string
}

type UpdateCategoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewUpdateCategoryUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, input UpdateCategoryInput, updateData dtos.UpdateCategoryInput) (*dtos.CategoryDTO, error) {
	category, err := uc.categoryRepo.FindByID(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	category.Update(updateData.Name, updateData.Description)
	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	var dto dtos.CategoryDTO
	dto.FromEntity(category)
	return &dto, nil
}

