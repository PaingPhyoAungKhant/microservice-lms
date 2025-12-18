package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCategories_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	listCategoriesUC := usecases.NewListCategoriesUseCase(categoryRepo, logger)

	ctx := context.Background()

	categories := []*entities.Category{
		entities.NewCategory("Category 1", "Description 1"),
		entities.NewCategory("Category 2", "Description 2"),
		entities.NewCategory("Category 3", "Description 3"),
	}

	for _, category := range categories {
		err := categoryRepo.Create(ctx, category)
		require.NoError(t, err)
	}

	limit := 10
	input := usecases.ListCategoriesInput{
		Limit: &limit,
	}

	result, err := listCategoriesUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.Total)
	assert.Len(t, result.Categories, 3)
}

func TestListCategories_Integration_WithSearch(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	listCategoriesUC := usecases.NewListCategoriesUseCase(categoryRepo, logger)

	ctx := context.Background()

	categories := []*entities.Category{
		entities.NewCategory("Programming", "Programming courses"),
		entities.NewCategory("Design", "Design courses"),
		entities.NewCategory("Marketing", "Marketing courses"),
	}

	for _, category := range categories {
		err := categoryRepo.Create(ctx, category)
		require.NoError(t, err)
	}

	searchQuery := "Programming"
	limit := 10
	input := usecases.ListCategoriesInput{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	result, err := listCategoriesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Categories, 1)
	assert.Equal(t, "Programming", result.Categories[0].Name)
}

func TestListCategories_Integration_WithPagination(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	listCategoriesUC := usecases.NewListCategoriesUseCase(categoryRepo, logger)

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		category := entities.NewCategory(fmt.Sprintf("Category %d", i), "Description")
		err := categoryRepo.Create(ctx, category)
		require.NoError(t, err)
	}

	limit := 2
	offset := 0
	input := usecases.ListCategoriesInput{
		Limit:  &limit,
		Offset: &offset,
	}

	result, err := listCategoriesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 5, result.Total)
	assert.Len(t, result.Categories, 2)

	offset = 2
	input.Offset = &offset

	result2, err := listCategoriesUC.Execute(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, 5, result2.Total)
	assert.Len(t, result2.Categories, 2)
}

func TestListCategories_Integration_WithSorting(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	listCategoriesUC := usecases.NewListCategoriesUseCase(categoryRepo, logger)

	ctx := context.Background()

	names := []string{"Zebra", "Alpha", "Beta"}
	for _, name := range names {
		category := entities.NewCategory(name, "Description")
		err := categoryRepo.Create(ctx, category)
		require.NoError(t, err)
	}

	sortColumn := "name"
	sortDirection := repositories.SortDirectionASC
	limit := 10
	input := usecases.ListCategoriesInput{
		SortColumn:    &sortColumn,
		SortDirection: &sortDirection,
		Limit:         &limit,
	}

	result, err := listCategoriesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.Total)
	assert.Len(t, result.Categories, 3)
}

