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

func TestAssignInstructor_Integration_Success(t *testing.T) {
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
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	assignInstructorUC := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	result, err := assignInstructorUC.Execute(ctx, offering.ID, input, instructorUsername)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, offering.ID, result.CourseOfferingID)
	assert.Equal(t, instructorID, result.InstructorID)
	assert.Equal(t, instructorUsername, result.InstructorUsername)

	createdInstructor, err := instructorRepo.FindByOfferingID(ctx, offering.ID)
	require.NoError(t, err)
	require.Len(t, createdInstructor, 1)
	assert.Equal(t, instructorID, createdInstructor[0].InstructorID)

	publisher.AssertExpectations(t)
}

func TestAssignInstructor_Integration_OfferingNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	instructorID := "instructor-123"
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	assignInstructorUC := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	nonExistentOfferingID := uuid.New().String()
	_, err = assignInstructorUC.Execute(ctx, nonExistentOfferingID, input, "username")
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseOfferingNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestAssignInstructor_Integration_DuplicateAssignment(t *testing.T) {
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
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	assignInstructorUC := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	result1, err := assignInstructorUC.Execute(ctx, offering.ID, input, instructorUsername)
	require.NoError(t, err)
	require.NotNil(t, result1)

	_, err = assignInstructorUC.Execute(ctx, offering.ID, input, instructorUsername)
	require.NoError(t, err)

	instructors, err := instructorRepo.FindByOfferingID(ctx, offering.ID)
	require.NoError(t, err)
	require.Len(t, instructors, 1)

	publisher.AssertExpectations(t)
}

