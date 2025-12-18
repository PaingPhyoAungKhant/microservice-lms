package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRemoveInstructorFromOffering_Success(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	instructorID := "instructor-123"
	instructorToRemove := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-123",
		CourseOfferingID:   offeringID,
		InstructorID:       instructorID,
		InstructorUsername: "instructor_user",
	}

	otherInstructor := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-456",
		CourseOfferingID:   offeringID,
		InstructorID:       "instructor-456",
		InstructorUsername: "other_instructor",
	}

	instructors := []*entities.CourseOfferingInstructor{otherInstructor, instructorToRemove}

	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return(instructors, nil).Once()
	instructorRepo.On("Delete", mock.Anything, instructorToRemove.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeInstructorRemovedFromOffering, mock.Anything).Return(nil).Once()

	uc := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	err := uc.Execute(context.Background(), offeringID, instructorID)
	require.NoError(t, err)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.InstructorRemovedFromOfferingEvent)
	require.True(t, ok)
	assert.Equal(t, instructorToRemove.ID, event.ID)
	assert.Equal(t, offeringID, event.CourseOfferingID)
	assert.Equal(t, instructorID, event.InstructorID)
	assert.Equal(t, instructorToRemove.InstructorUsername, event.InstructorUsername)

	instructorRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestRemoveInstructorFromOffering_InstructorNotFound(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	instructorID := "instructor-123"
	otherInstructor := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-456",
		CourseOfferingID:   offeringID,
		InstructorID:       "instructor-456",
		InstructorUsername: "other_instructor",
	}

	instructors := []*entities.CourseOfferingInstructor{otherInstructor}

	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return(instructors, nil).Once()

	uc := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	err := uc.Execute(context.Background(), offeringID, instructorID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instructor not found")

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	instructorRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	instructorRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestRemoveInstructorFromOffering_EmptyInstructors(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	instructorID := "instructor-123"

	instructors := []*entities.CourseOfferingInstructor{}

	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return(instructors, nil).Once()

	uc := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	err := uc.Execute(context.Background(), offeringID, instructorID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instructor not found")

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	instructorRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	instructorRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestRemoveInstructorFromOffering_RepositoryError(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	instructorID := "instructor-123"
	instructorToRemove := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-123",
		CourseOfferingID:   offeringID,
		InstructorID:       instructorID,
		InstructorUsername: "instructor_user",
	}

	instructors := []*entities.CourseOfferingInstructor{instructorToRemove}

	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return(instructors, nil).Once()
	instructorRepo.On("Delete", mock.Anything, instructorToRemove.ID).Return(assert.AnError).Once()

	uc := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, publisher, logger)

	err := uc.Execute(context.Background(), offeringID, instructorID)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	instructorRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

