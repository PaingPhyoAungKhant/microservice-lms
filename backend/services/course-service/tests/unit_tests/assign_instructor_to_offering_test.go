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

func TestAssignInstructorToOffering_Success(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     "course-123",
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	instructorID := "instructor-123"
	instructorUsername := "instructor_user"
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	var createdInstructor *entities.CourseOfferingInstructor
	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(offering, nil).Once()
	instructorRepo.On("Create", mock.Anything, mock.MatchedBy(func(i *entities.CourseOfferingInstructor) bool {
		createdInstructor = i
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeInstructorAssignedToOffering, mock.Anything).Return(nil).Once()

	uc := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), offeringID, input, instructorUsername)
	require.NoError(t, err)
	require.NotNil(t, createdInstructor)
	require.Equal(t, offeringID, createdInstructor.CourseOfferingID)
	require.Equal(t, instructorID, createdInstructor.InstructorID)
	require.Equal(t, instructorUsername, createdInstructor.InstructorUsername)

	assertDTOEqualCourseOfferingInstructor(t, dto, createdInstructor)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.InstructorAssignedToOfferingEvent)
	require.True(t, ok)
	assert.Equal(t, createdInstructor.ID, event.ID)
	assert.Equal(t, offeringID, event.CourseOfferingID)
	assert.Equal(t, instructorID, event.InstructorID)
	assert.Equal(t, instructorUsername, event.InstructorUsername)

	instructorRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestAssignInstructorToOffering_OfferingNotFound(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	instructorID := "instructor-123"
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return((*entities.CourseOffering)(nil), nil).Once()

	uc := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input, "username")
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseOfferingNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	instructorRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	instructorRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestAssignInstructorToOffering_RepositoryError(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     "course-123",
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	instructorID := "instructor-123"
	input := dtos.AssignInstructorInput{
		InstructorID: instructorID,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(offering, nil).Once()
	instructorRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input, "username")
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	instructorRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

