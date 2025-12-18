package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEnrollment_AlreadyExists(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewCreateEnrollmentUseCase(repo, publisher, logger)

	studentID := "student-id-123"
	courseOfferingID := "offering-id-123"
	limit := 1

	existingEnrollment := entities.NewEnrollment(
		studentID,
		"existingstudent",
		"course-id-123",
		"Existing Course",
		courseOfferingID,
		"Fall 2024",
	)

	repo.On("Find", mock.Anything, mock.MatchedBy(func(q repositories.EnrollmentQuery) bool {
		return q.StudentID != nil && *q.StudentID == studentID &&
			q.CourseOfferingID != nil && *q.CourseOfferingID == courseOfferingID &&
			q.Limit != nil && *q.Limit == limit
	})).Return(&repositories.EnrollmentQueryResult{
		Enrollments: []*entities.Enrollment{existingEnrollment},
		Total:      1,
	}, nil).Once()

	_, err := uc.Execute(context.Background(), usecases.CreateEnrollmentInput{
		StudentID:          studentID,
		StudentUsername:    "existingstudent",
		CourseID:           "course-id-123",
		CourseName:         "Existing Course",
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: "Fall 2024",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentAlreadyExists, err)

	repo.AssertExpectations(t)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateEnrollment_Success(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewCreateEnrollmentUseCase(repo, publisher, logger)

	studentID := "student-id-new"
	courseOfferingID := "offering-id-new"
	limit := 1

	repo.On("Find", mock.Anything, mock.MatchedBy(func(q repositories.EnrollmentQuery) bool {
		return q.StudentID != nil && *q.StudentID == studentID &&
			q.CourseOfferingID != nil && *q.CourseOfferingID == courseOfferingID &&
			q.Limit != nil && *q.Limit == limit
	})).Return(&repositories.EnrollmentQueryResult{
		Enrollments: []*entities.Enrollment{},
		Total:      0,
	}, nil).Once()

	var createdEnrollment *entities.Enrollment
	repo.On("Create", mock.Anything, mock.MatchedBy(func(e *entities.Enrollment) bool {
		createdEnrollment = e
		return e.StudentID == studentID
	})).Return(nil).Once()

	publisher.On("Publish", mock.Anything, events.EventTypeEnrollmentCreated, mock.Anything).Return(nil).Once()

	result, err := uc.Execute(context.Background(), usecases.CreateEnrollmentInput{
		StudentID:          studentID,
		StudentUsername:    "newstudent",
		CourseID:           "course-id-new",
		CourseName:         "New Course",
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: "Fall 2024",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, createdEnrollment)
	assert.Equal(t, studentID, result.StudentID)
	assert.Equal(t, studentID, createdEnrollment.StudentID)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateEnrollment_RepositoryError(t *testing.T) {
	repo := new(mocks.MockEnrollmentRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewCreateEnrollmentUseCase(repo, publisher, logger)

	studentID := "student-id-error"
	courseOfferingID := "offering-id-error"

	repo.On("Find", mock.Anything, mock.Anything).Return(&repositories.EnrollmentQueryResult{
		Enrollments: []*entities.Enrollment{},
		Total:      0,
	}, nil).Once()

	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()

	_, err := uc.Execute(context.Background(), usecases.CreateEnrollmentInput{
		StudentID:          studentID,
		StudentUsername:    "errorstudent",
		CourseID:           "course-id-error",
		CourseName:         "Error Course",
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: "Fall 2024",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	repo.AssertExpectations(t)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

