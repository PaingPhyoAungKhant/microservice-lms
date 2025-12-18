package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	createCategoryUC := usecases.NewCreateCategoryUseCase(categoryRepo, logger)

	input := dtos.CreateCategoryInput{
		Name:        "Programming",
		Description: "Programming related courses",
	}

	ctx := context.Background()
	result, err := createCategoryUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)

	createdCategory, err := categoryRepo.FindByName(ctx, input.Name)
	require.NoError(t, err)
	require.NotNil(t, createdCategory)
	assert.Equal(t, input.Name, createdCategory.Name)
	assert.Equal(t, input.Description, createdCategory.Description)
}

func TestCreateCategory_Integration_DuplicateName(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	createCategoryUC := usecases.NewCreateCategoryUseCase(categoryRepo, logger)

	ctx := context.Background()

	input1 := dtos.CreateCategoryInput{
		Name:        "Programming",
		Description: "Programming related courses",
	}

	_, err = createCategoryUC.Execute(ctx, input1)
	require.NoError(t, err)

	input2 := dtos.CreateCategoryInput{
		Name:        "Programming",
		Description: "Different description",
	}

	_, err = createCategoryUC.Execute(ctx, input2)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCategoryAlreadyExists, err)
}

func TestCreateCategory_Integration_EmptyName(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"category"},
	})
	require.NoError(t, err)
	defer cleanup()

	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()

	createCategoryUC := usecases.NewCreateCategoryUseCase(categoryRepo, logger)

	input := dtos.CreateCategoryInput{
		Name:        "",
		Description: "Description",
	}

	ctx := context.Background()
	_, err = createCategoryUC.Execute(ctx, input)

	require.Error(t, err)
}

