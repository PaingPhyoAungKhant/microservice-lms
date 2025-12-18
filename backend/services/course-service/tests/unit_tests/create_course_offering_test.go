package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourseOffering_Success(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	courseRepo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	duration := "3 months"
	classTime := "Mon, Wed, Fri 10:00 AM"
	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "online",
		Duration:       &duration,
		ClassTime:      &classTime,
		EnrollmentCost: 99.99,
	}

	var createdOffering *entities.CourseOffering
	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("Create", mock.Anything, mock.MatchedBy(func(o *entities.CourseOffering) bool {
		createdOffering = o
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseOfferingCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), courseID, input)
	require.NoError(t, err)
	require.NotNil(t, createdOffering)
	require.Equal(t, input.Name, createdOffering.Name)
	require.Equal(t, input.Description, createdOffering.Description)
	require.Equal(t, entities.OfferingTypeOnline, createdOffering.OfferingType)
	require.Equal(t, entities.OfferingStatusPending, createdOffering.Status)
	require.Equal(t, duration, *createdOffering.Duration)
	require.Equal(t, classTime, *createdOffering.ClassTime)
	require.Equal(t, input.EnrollmentCost, createdOffering.EnrollmentCost)

	assertDTOEqualCourseOffering(t, dto, createdOffering)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseOfferingCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdOffering.ID, event.ID)
	assert.Equal(t, courseID, event.CourseID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, string(createdOffering.OfferingType), event.OfferingType)
	assert.Equal(t, input.EnrollmentCost, event.EnrollmentCost)

	offeringRepo.AssertExpectations(t)
	courseRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourseOffering_WithoutDurationAndClassTime(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	courseRepo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "oncampus",
		EnrollmentCost: 0.0,
	}

	var createdOffering *entities.CourseOffering
	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("Create", mock.Anything, mock.MatchedBy(func(o *entities.CourseOffering) bool {
		createdOffering = o
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseOfferingCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), courseID, input)
	require.NoError(t, err)
	require.NotNil(t, createdOffering)
	require.Nil(t, createdOffering.Duration)
	require.Nil(t, createdOffering.ClassTime)
	require.Equal(t, entities.OfferingTypeOnCampus, createdOffering.OfferingType)

	assertDTOEqualCourseOffering(t, dto, createdOffering)

	offeringRepo.AssertExpectations(t)
	courseRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourseOffering_CourseNotFound(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	courseRepo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	courseID := "course-123"
	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	courseRepo.On("FindByID", mock.Anything, courseID).Return((*entities.Course)(nil), nil).Once()

	uc := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), courseID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	offeringRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	offeringRepo.AssertExpectations(t)
	courseRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourseOffering_RepositoryError(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	courseRepo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	input := dtos.CreateCourseOfferingInput{
		Name:           "Spring 2024",
		Description:    "Spring semester offering",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), courseID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	offeringRepo.AssertExpectations(t)
	courseRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

