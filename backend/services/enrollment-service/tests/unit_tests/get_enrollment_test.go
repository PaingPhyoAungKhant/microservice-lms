package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetEnrollment_Success(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	logger := logger.NewNop()

	uc := usecases.NewGetEnrollmentUseCase(repo, logger)

	enrollment := entities.NewEnrollment(
		"student-id-123",
		"teststudent",
		"course-id-123",
		"Test Course",
		"offering-id-123",
		"Fall 2024",
	)

	repo.On("FindByID", mock.Anything, enrollment.ID).Return(enrollment, nil).Once()

	result, err := uc.Execute(context.Background(), usecases.GetEnrollmentInput{
		EnrollmentID: enrollment.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, enrollment.ID, result.ID)
	assert.Equal(t, enrollment.StudentID, result.StudentID)

	repo.AssertExpectations(t)
}

func TestGetEnrollment_NotFound(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	logger := logger.NewNop()

	uc := usecases.NewGetEnrollmentUseCase(repo, logger)

	repo.On("FindByID", mock.Anything, "non-existent-id").Return((*entities.Enrollment)(nil), errors.New("enrollment not found")).Once()

	_, err := uc.Execute(context.Background(), usecases.GetEnrollmentInput{
		EnrollmentID: "non-existent-id",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentNotFound, err)

	repo.AssertExpectations(t)
}

func TestGetEnrollment_NilEnrollment(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	logger := logger.NewNop()

	uc := usecases.NewGetEnrollmentUseCase(repo, logger)

	repo.On("FindByID", mock.Anything, "nil-id").Return((*entities.Enrollment)(nil), nil).Once()

	_, err := uc.Execute(context.Background(), usecases.GetEnrollmentInput{
		EnrollmentID: "nil-id",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentNotFound, err)

	repo.AssertExpectations(t)
}

