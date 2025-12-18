package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

var (
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
	ErrCategoryNotFound      = errors.New("category not found")
)

type CreateCategoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewCreateCategoryUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input dtos.CreateCategoryInput) (*dtos.CategoryDTO, error) {
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	existing, _ := uc.categoryRepo.FindByName(ctx, input.Name)
	if existing != nil {
		return nil, ErrCategoryAlreadyExists
	}

	category := entities.NewCategory(input.Name, input.Description)
	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		uc.logger.Error("failed to create category", zap.Error(err))
		return nil, err
	}

	var dto dtos.CategoryDTO
	dto.FromEntity(category)
	return &dto, nil
}

