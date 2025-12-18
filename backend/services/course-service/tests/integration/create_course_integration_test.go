package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourse_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	courseCategoryRepo := SetupCourseCategoryRepository(db)
	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"
	publisher := new(sharedMocks.MockPublisher)

	createCourseUC := usecases.NewCreateCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	thumbnailID := uuid.New().String()
	input := dtos.CreateCourseInput{
		Name:        "Introduction to Go",
		Description: "Learn Go programming language",
		ThumbnailID: &thumbnailID,
	}

	ctx := context.Background()
	result, err := createCourseUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, thumbnailID, *result.ThumbnailID)
	assert.Contains(t, result.ThumbnailURL, apiGatewayURL)
	assert.Contains(t, result.ThumbnailURL, "/api/v1/buckets/course-thumbnails/files/")
	assert.Contains(t, result.ThumbnailURL, thumbnailID)

	createdCourse, err := courseRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, createdCourse)
	assert.Equal(t, input.Name, createdCourse.Name)
	assert.Equal(t, input.Description, createdCourse.Description)

	publisher.AssertExpectations(t)
}

func TestCreateCourse_Integration_WithoutThumbnail(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	courseCategoryRepo := SetupCourseCategoryRepository(db)
	categoryRepo := SetupCategoryRepository(db)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"
	publisher := new(sharedMocks.MockPublisher)

	createCourseUC := usecases.NewCreateCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	input := dtos.CreateCourseInput{
		Name:        "Introduction to Python",
		Description: "Learn Python programming",
	}

	ctx := context.Background()
	result, err := createCourseUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Nil(t, result.ThumbnailID)
	assert.Empty(t, result.ThumbnailURL)

	createdCourse, err := courseRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, createdCourse)
	assert.Nil(t, createdCourse.ThumbnailID)

	publisher.AssertExpectations(t)
}

