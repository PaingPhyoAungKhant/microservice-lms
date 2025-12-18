package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRemoveInstructor_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	instructorID := uuid.New().String()
	instructorUsername := "instructor_user"
	instructor := entities.NewCourseOfferingInstructor(offering.ID, instructorID, instructorUsername)
	err = instructorRepo.Create(ctx, instructor)
	require.NoError(t, err)

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	removeInstructorUC := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	err = removeInstructorUC.Execute(ctx, offering.ID, instructorID)
	require.NoError(t, err)

	instructors, err := instructorRepo.FindByOfferingID(ctx, offering.ID)
	require.NoError(t, err)
	require.Len(t, instructors, 0)

	publisher.AssertExpectations(t)
}

func TestRemoveInstructor_Integration_InstructorNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	removeInstructorUC := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	nonExistentInstructorID := uuid.New().String()
	err = removeInstructorUC.Execute(ctx, offering.ID, nonExistentInstructorID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instructor not found")

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

