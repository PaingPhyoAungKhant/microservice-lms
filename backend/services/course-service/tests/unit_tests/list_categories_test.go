package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCategories_Success(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()

	category1 := entities.NewCategory("Category 1", "Description 1")
	category2 := entities.NewCategory("Category 2", "Description 2")

	limit := 10
	query := repositories.CategoryQuery{
		Limit: &limit,
	}

	result := &repositories.CategoryQueryResult{
		Categories: []*entities.Category{category1, category2},
		Total:      2,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()

	uc := usecases.NewListCategoriesUseCase(repo, logger)

	input := usecases.ListCategoriesInput{
		Limit: &limit,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.Total)
	assert.Len(t, output.Categories, 2)

	repo.AssertExpectations(t)
}

func TestListCategories_WithSearch(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()

	searchQuery := "test"
	limit := 10
	query := repositories.CategoryQuery{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	category := entities.NewCategory("Test Category", "Description")
	result := &repositories.CategoryQueryResult{
		Categories: []*entities.Category{category},
		Total:      1,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()

	uc := usecases.NewListCategoriesUseCase(repo, logger)

	input := usecases.ListCategoriesInput{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 1, output.Total)
	assert.Len(t, output.Categories, 1)

	repo.AssertExpectations(t)
}

func TestListCategories_WithPagination(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()

	limit := 2
	offset := 0
	query := repositories.CategoryQuery{
		Limit:  &limit,
		Offset: &offset,
	}

	category1 := entities.NewCategory("Category 1", "Description 1")
	category2 := entities.NewCategory("Category 2", "Description 2")

	result := &repositories.CategoryQueryResult{
		Categories: []*entities.Category{category1, category2},
		Total:      5,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()

	uc := usecases.NewListCategoriesUseCase(repo, logger)

	input := usecases.ListCategoriesInput{
		Limit:  &limit,
		Offset: &offset,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 5, output.Total)
	assert.Len(t, output.Categories, 2)

	repo.AssertExpectations(t)
}

