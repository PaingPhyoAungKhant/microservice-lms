package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteEnrollment_Success(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewDeleteEnrollmentUseCase(repo, publisher, logger)

	enrollment := entities.NewEnrollment(
		"student-id-123",
		"teststudent",
		"course-id-123",
		"Test Course",
		"offering-id-123",
		"Fall 2024",
	)

	repo.On("FindByID", mock.Anything, enrollment.ID).Return(enrollment, nil).Once()
	repo.On("Delete", mock.Anything, enrollment.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeEnrollmentDeleted, mock.Anything).Return(nil).Once()

	result, err := uc.Execute(context.Background(), usecases.DeleteEnrollmentInput{
		EnrollmentID: enrollment.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Message, "deleted successfully")

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteEnrollment_NotFound(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewDeleteEnrollmentUseCase(repo, publisher, logger)

	repo.On("FindByID", mock.Anything, "non-existent-id").Return((*entities.Enrollment)(nil), errors.New("enrollment not found")).Once()

	_, err := uc.Execute(context.Background(), usecases.DeleteEnrollmentInput{
		EnrollmentID: "non-existent-id",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentNotFound, err)

	repo.AssertExpectations(t)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteEnrollment_DeleteError(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewDeleteEnrollmentUseCase(repo, publisher, logger)

	enrollment := entities.NewEnrollment(
		"student-id-error",
		"errorstudent",
		"course-id-error",
		"Error Course",
		"offering-id-error",
		"Fall 2024",
	)

	repo.On("FindByID", mock.Anything, enrollment.ID).Return(enrollment, nil).Once()
	repo.On("Delete", mock.Anything, enrollment.ID).Return(errors.New("database error")).Once()

	_, err := uc.Execute(context.Background(), usecases.DeleteEnrollmentInput{
		EnrollmentID: enrollment.ID,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete enrollment")

	repo.AssertExpectations(t)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

