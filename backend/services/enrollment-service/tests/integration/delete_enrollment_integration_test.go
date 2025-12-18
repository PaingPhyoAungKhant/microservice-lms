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

func TestDeleteEnrollment_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	enrollment := entities.NewEnrollment(
		uuid.New().String(),
		"deletestudent",
		uuid.New().String(),
		"Delete Course",
		uuid.New().String(),
		"Fall 2024",
	)
	enrollmentRepo.Create(ctx, enrollment)

	deleteUC := usecases.NewDeleteEnrollmentUseCase(enrollmentRepo, publisher, logger)

	input := usecases.DeleteEnrollmentInput{
		EnrollmentID: enrollment.ID,
	}

	publisher.On("Publish", mock.Anything, events.EventTypeEnrollmentDeleted, mock.Anything).Return(nil).Once()

	result, err := deleteUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Message, "deleted successfully")

	_, err = enrollmentRepo.FindByID(ctx, enrollment.ID)
	require.Error(t, err)

	publisher.AssertExpectations(t)
}

