package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)


type GetCategoryInput struct {
	CategoryID string
}

type GetCategoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewGetCategoryUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *GetCategoryUseCase {
	return &GetCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *GetCategoryUseCase) Execute(ctx context.Context, input GetCategoryInput) (*dtos.CategoryDTO, error) {
	category, err := uc.categoryRepo.FindByID(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	var dto dtos.CategoryDTO
	dto.FromEntity(category)
	return &dto, nil
}

