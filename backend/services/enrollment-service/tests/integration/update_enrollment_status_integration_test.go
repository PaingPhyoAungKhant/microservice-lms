package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateEnrollmentStatus_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	enrollment := entities.NewEnrollment(
		uuid.New().String(),
		"updatestudent",
		uuid.New().String(),
		"Update Course",
		uuid.New().String(),
		"Fall 2024",
	)
	enrollmentRepo.Create(ctx, enrollment)

	status, _ := valueobjects.NewEnrollmentStatus("approved")
	updateUC := usecases.NewUpdateEnrollmentStatusUseCase(enrollmentRepo, publisher, logger)

	input := usecases.UpdateEnrollmentStatusInput{
		EnrollmentID: enrollment.ID,
		Status:       status,
	}

	publisher.On("Publish", mock.Anything, events.EventTypeEnrollmentUpdated, mock.Anything).Return(nil).Once()

	result, err := updateUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "approved", result.Status)

	updatedEnrollment, err := enrollmentRepo.FindByID(ctx, enrollment.ID)
	require.NoError(t, err)
	assert.Equal(t, "approved", updatedEnrollment.Status.String())

	publisher.AssertExpectations(t)
}

func TestUpdateEnrollmentStatus_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	status, _ := valueobjects.NewEnrollmentStatus("approved")
	updateUC := usecases.NewUpdateEnrollmentStatusUseCase(enrollmentRepo, publisher, logger)

	input := usecases.UpdateEnrollmentStatusInput{
		EnrollmentID: "00000000-0000-0000-0000-000000000000",
		Status:       status,
	}

	_, err := updateUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

