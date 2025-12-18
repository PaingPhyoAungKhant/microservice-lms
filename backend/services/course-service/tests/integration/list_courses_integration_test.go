package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCourses_Integration_Success(t *testing.T) {
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

	listCoursesUC := usecases.NewListCoursesUseCase(courseRepo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	ctx := context.Background()

	courses := []*entities.Course{
		entities.NewCourse("Course 1", "Description 1", nil),
		entities.NewCourse("Course 2", "Description 2", nil),
		entities.NewCourse("Course 3", "Description 3", nil),
	}

	for _, course := range courses {
		err := courseRepo.Create(ctx, course)
		require.NoError(t, err)
	}

	limit := 10
	input := usecases.ListCoursesInput{
		Limit: &limit,
	}

	result, err := listCoursesUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.Total)
	assert.Len(t, result.Courses, 3)
}

func TestListCourses_Integration_WithSearch(t *testing.T) {
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

	listCoursesUC := usecases.NewListCoursesUseCase(courseRepo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	ctx := context.Background()

	courses := []*entities.Course{
		entities.NewCourse("Introduction to Go", "Learn Go programming", nil),
		entities.NewCourse("Introduction to Python", "Learn Python programming", nil),
		entities.NewCourse("Web Design Basics", "Learn web design", nil),
	}

	for _, course := range courses {
		err := courseRepo.Create(ctx, course)
		require.NoError(t, err)
	}

	searchQuery := "Go"
	limit := 10
	input := usecases.ListCoursesInput{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	result, err := listCoursesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Courses, 1)
	assert.Contains(t, result.Courses[0].Name, "Go")
}

func TestListCourses_Integration_WithPagination(t *testing.T) {
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

	listCoursesUC := usecases.NewListCoursesUseCase(courseRepo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		course := entities.NewCourse(fmt.Sprintf("Course %d", i), "Description", nil)
		err := courseRepo.Create(ctx, course)
		require.NoError(t, err)
	}

	limit := 2
	offset := 0
	input := usecases.ListCoursesInput{
		Limit:  &limit,
		Offset: &offset,
	}

	result, err := listCoursesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 5, result.Total)
	assert.Len(t, result.Courses, 2)

	offset = 2
	input.Offset = &offset

	result2, err := listCoursesUC.Execute(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, 5, result2.Total)
	assert.Len(t, result2.Courses, 2)
}

func TestListCourses_Integration_WithThumbnailURL(t *testing.T) {
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

	listCoursesUC := usecases.NewListCoursesUseCase(courseRepo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	ctx := context.Background()

	thumbnailID := uuid.New().String()
	course := entities.NewCourse("Course with Thumbnail", "Description", &thumbnailID)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	limit := 10
	input := usecases.ListCoursesInput{
		Limit: &limit,
	}

	result, err := listCoursesUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Courses, 1)
	assert.NotEmpty(t, result.Courses[0].ThumbnailURL)
	assert.Contains(t, result.Courses[0].ThumbnailURL, apiGatewayURL)
	assert.Contains(t, result.Courses[0].ThumbnailURL, "/api/v1/buckets/course-thumbnails/files/")
	assert.Contains(t, result.Courses[0].ThumbnailURL, thumbnailID)
}

