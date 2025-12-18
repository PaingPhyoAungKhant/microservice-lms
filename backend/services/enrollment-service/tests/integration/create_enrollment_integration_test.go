package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEnrollment_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	createUC := usecases.NewCreateEnrollmentUseCase(enrollmentRepo, publisher, logger)

	studentID := uuid.New().String()
	courseID := uuid.New().String()
	courseOfferingID := uuid.New().String()

	input := usecases.CreateEnrollmentInput{
		StudentID:          studentID,
		StudentUsername:    "integrationstudent",
		CourseID:           courseID,
		CourseName:         "Integration Course",
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: "Fall 2024",
	}

	publisher.On("Publish", mock.Anything, events.EventTypeEnrollmentCreated, mock.Anything).Return(nil).Once()

	result, err := createUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.StudentID, result.StudentID)
	assert.Equal(t, input.StudentUsername, result.StudentUsername)
	assert.Equal(t, input.CourseID, result.CourseID)
	assert.NotEmpty(t, result.ID)

	enrollment, err := enrollmentRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, result.ID, enrollment.ID)

	publisher.AssertExpectations(t)
}

func TestCreateEnrollment_Integration_DuplicateEnrollment(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	studentID := uuid.New().String()
	courseID := uuid.New().String()
	courseOfferingID := uuid.New().String()

	enrollment := entities.NewEnrollment(
		studentID,
		"duplicatestudent",
		courseID,
		"Duplicate Course",
		courseOfferingID,
		"Fall 2024",
	)
	enrollmentRepo.Create(ctx, enrollment)

	createUC := usecases.NewCreateEnrollmentUseCase(enrollmentRepo, publisher, logger)

	input := usecases.CreateEnrollmentInput{
		StudentID:          studentID,
		StudentUsername:    "duplicatestudent",
		CourseID:           courseID,
		CourseName:         "Duplicate Course",
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: "Fall 2024",
	}

	_, err := createUC.Execute(ctx, input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentAlreadyExists, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

