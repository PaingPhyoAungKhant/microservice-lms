package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFindEnrollment_Success(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	logger := logger.NewNop()

	uc := usecases.NewFindEnrollmentUseCase(repo, logger)

	enrollment1 := entities.NewEnrollment(
		"student-id-1",
		"student1",
		"course-id-1",
		"Course 1",
		"offering-id-1",
		"Fall 2024",
	)
	enrollment2 := entities.NewEnrollment(
		"student-id-2",
		"student2",
		"course-id-2",
		"Course 2",
		"offering-id-2",
		"Spring 2024",
	)

	limit := 10
	repo.On("Find", mock.Anything, mock.MatchedBy(func(q repositories.EnrollmentQuery) bool {
		return q.Limit != nil && *q.Limit == limit
	})).Return(&repositories.EnrollmentQueryResult{
		Enrollments: []*entities.Enrollment{enrollment1, enrollment2},
		Total:      2,
	}, nil).Once()

	result, err := uc.Execute(context.Background(), usecases.FindEnrollmentInput{
		Limit: &limit,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Enrollments, 2)

	repo.AssertExpectations(t)
}

func TestFindEnrollment_WithFilters(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	logger := logger.NewNop()

	uc := usecases.NewFindEnrollmentUseCase(repo, logger)

	enrollment := entities.NewEnrollment(
		"student-id-filter",
		"filterstudent",
		"course-id-filter",
		"Filter Course",
		"offering-id-filter",
		"Fall 2024",
	)

	status, _ := valueobjects.NewEnrollmentStatus("pending")
	studentID := enrollment.StudentID
	limit := 10

	repo.On("Find", mock.Anything, mock.MatchedBy(func(q repositories.EnrollmentQuery) bool {
		return q.StudentID != nil && *q.StudentID == studentID &&
			q.Status != nil && q.Status.String() == "pending" &&
			q.Limit != nil && *q.Limit == limit
	})).Return(&repositories.EnrollmentQueryResult{
		Enrollments: []*entities.Enrollment{enrollment},
		Total:      1,
	}, nil).Once()

	result, err := uc.Execute(context.Background(), usecases.FindEnrollmentInput{
		StudentID: &studentID,
		Status:   &status,
		Limit:    &limit,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Enrollments, 1)

	repo.AssertExpectations(t)
}

