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

func TestUpdateCourseOffering_Integration_Success(t *testing.T) {
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

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 99.99)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	duration := "6 months"
	classTime := "Tue, Thu 2:00 PM"
	input := dtos.UpdateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Updated description",
		OfferingType:   "oncampus",
		Duration:       &duration,
		ClassTime:      &classTime,
		EnrollmentCost: 149.99,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	updateOfferingUC := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, publisher, logger)

	result, err := updateOfferingUC.Execute(ctx, offering.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, "oncampus", result.OfferingType)
	assert.Equal(t, duration, *result.Duration)
	assert.Equal(t, classTime, *result.ClassTime)
	assert.Equal(t, input.EnrollmentCost, result.EnrollmentCost)

	updatedOffering, err := offeringRepo.FindByID(ctx, offering.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedOffering)
	assert.Equal(t, input.Name, updatedOffering.Name)
	assert.Equal(t, input.Description, updatedOffering.Description)

	publisher.AssertExpectations(t)
}

func TestUpdateCourseOffering_Integration_OfferingNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	offeringRepo := SetupCourseOfferingRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.UpdateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Updated description",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	updateOfferingUC := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, publisher, logger)

	nonExistentOfferingID := uuid.New().String()
	_, err = updateOfferingUC.Execute(ctx, nonExistentOfferingID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseOfferingNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

