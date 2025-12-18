package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type ListCategoriesInput struct {
	SearchQuery   *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *repositories.SortDirection
}

type ListCategoriesOutput struct {
	Categories []dtos.CategoryDTO `json:"categories"`
	Total      int                 `json:"total"`
}

type ListCategoriesUseCase struct {
	categoryRepo repositories.CategoryRepository
	logger       *logger.Logger
}

func NewListCategoriesUseCase(categoryRepo repositories.CategoryRepository, logger *logger.Logger) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *ListCategoriesUseCase) Execute(ctx context.Context, input ListCategoriesInput) (*ListCategoriesOutput, error) {
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

	return &ListCategoriesOutput{
		Categories: categoryDTOs,
		Total:      result.Total,
	}, nil
}

