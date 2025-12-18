package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindCategoryInput struct {
	SearchQuery   *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *repositories.SortDirection
}

type FindCategoryOutput struct {
	Categories []dtos.CategoryDTO `json:"categories"`
	Total      int                 `json:"total"`
}

type FindCategoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewFindCategoryUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *FindCategoryUseCase {
	return &FindCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *FindCategoryUseCase) Execute(ctx context.Context, input FindCategoryInput) (*FindCategoryOutput, error) {
	query := repositories.CategoryQuery{
		SearchQuery:   input.SearchQuery,
		Limit:         input.Limit,
		Offset:        input.Offset,
		SortColumn:    input.SortColumn,
		SortDirection: input.SortDirection,
	}

	result, err := uc.categoryRepo.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	categoryDTOs := make([]dtos.CategoryDTO, len(result.Categories))
	for i, category := range result.Categories {
		categoryDTOs[i].FromEntity(category)
	}

	return &FindCategoryOutput{
		Categories: categoryDTOs,
		Total:      result.Total,
	}, nil
}

