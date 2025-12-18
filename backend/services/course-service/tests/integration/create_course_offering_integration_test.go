package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourseOffering_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	duration := "3 months"
	classTime := "Mon, Wed, Fri 10:00 AM"
	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "online",
		Duration:       &duration,
		ClassTime:      &classTime,
		EnrollmentCost: 99.99,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	createOfferingUC := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	result, err := createOfferingUC.Execute(ctx, course.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, "online", result.OfferingType)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, duration, *result.Duration)
	assert.Equal(t, classTime, *result.ClassTime)
	assert.Equal(t, input.EnrollmentCost, result.EnrollmentCost)

	createdOffering, err := offeringRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, createdOffering)
	assert.Equal(t, input.Name, createdOffering.Name)
	assert.Equal(t, course.ID, createdOffering.CourseID)

	publisher.AssertExpectations(t)
}

func TestCreateCourseOffering_Integration_CourseNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	createOfferingUC := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	nonExistentCourseID := uuid.New().String()
	_, err = createOfferingUC.Execute(ctx, nonExistentCourseID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateCourseOffering_Integration_WithoutDurationAndClassTime(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	input := dtos.CreateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Fall semester offering",
		OfferingType:   "oncampus",
		EnrollmentCost: 0.0,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	createOfferingUC := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	result, err := createOfferingUC.Execute(ctx, course.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Duration)
	assert.Nil(t, result.ClassTime)
	assert.Equal(t, "oncampus", result.OfferingType)

	publisher.AssertExpectations(t)
}

