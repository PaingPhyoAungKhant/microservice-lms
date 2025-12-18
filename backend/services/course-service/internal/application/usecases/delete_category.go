package usecases

import (
	"context"
	"fmt"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)


type DeleteCategoryInput struct {
	CategoryID string
}

type DeleteCategoryOutput struct {
	Message string
}

type DeleteCategoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewDeleteCategoryUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, input DeleteCategoryInput) (*DeleteCategoryOutput, error) {
	category, err := uc.categoryRepo.FindByID(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if err := uc.categoryRepo.Delete(ctx, input.CategoryID); err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return &DeleteCategoryOutput{
		Message: "Category deleted successfully",
	}, nil
}

