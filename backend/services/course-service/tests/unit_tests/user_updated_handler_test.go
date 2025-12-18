package unit_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserUpdatedHandler_Success(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	logger := logger.NewNop()

	userID := "user-123"
	newUsername := "new_username"
	oldUsername := "old_username"

	instructor1 := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-1",
		CourseOfferingID:   "offering-1",
		InstructorID:       userID,
		InstructorUsername: oldUsername,
	}

	instructor2 := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-2",
		CourseOfferingID:   "offering-2",
		InstructorID:       userID,
		InstructorUsername: oldUsername,
	}

	instructors := []*entities.CourseOfferingInstructor{instructor1, instructor2}

	event := events.UserUpdatedEvent{
		ID:            userID,
		Email:         "user@example.com",
		Username:      newUsername,
		Role:          "instructor",
		Status:        "active",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	eventBody, err := json.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	instructorRepo.On("FindByInstructorID", ctx, userID).Return(instructors, nil).Once()
	instructorRepo.On("Update", ctx, mock.MatchedBy(func(i *entities.CourseOfferingInstructor) bool {
		return i.ID == instructor1.ID && i.InstructorUsername == newUsername
	})).Return(nil).Once()
	instructorRepo.On("Update", ctx, mock.MatchedBy(func(i *entities.CourseOfferingInstructor) bool {
		return i.ID == instructor2.ID && i.InstructorUsername == newUsername
	})).Return(nil).Once()

	handler := handlers.NewUserUpdatedHandler(instructorRepo, logger)

	err = handler.Handle(eventBody)
	require.NoError(t, err)

	require.Equal(t, newUsername, instructor1.InstructorUsername)
	require.Equal(t, newUsername, instructor2.InstructorUsername)

	instructorRepo.AssertExpectations(t)
}

func TestUserUpdatedHandler_NoInstructorsFound(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	logger := logger.NewNop()

	userID := "user-123"
	newUsername := "new_username"

	event := events.UserUpdatedEvent{
		ID:            userID,
		Email:         "user@example.com",
		Username:      newUsername,
		Role:          "instructor",
		Status:        "active",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	eventBody, err := json.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	instructors := []*entities.CourseOfferingInstructor{}
	instructorRepo.On("FindByInstructorID", ctx, userID).Return(instructors, nil).Once()

	handler := handlers.NewUserUpdatedHandler(instructorRepo, logger)

	err = handler.Handle(eventBody)
	require.NoError(t, err)

	instructorRepo.AssertNotCalled(t, "Update", ctx, mock.Anything)
	instructorRepo.AssertExpectations(t)
}

func TestUserUpdatedHandler_UsernameUnchanged(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	logger := logger.NewNop()

	userID := "user-123"
	username := "same_username"

	instructor := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-1",
		CourseOfferingID:   "offering-1",
		InstructorID:       userID,
		InstructorUsername: username,
	}

	instructors := []*entities.CourseOfferingInstructor{instructor}

	event := events.UserUpdatedEvent{
		ID:            userID,
		Email:         "user@example.com",
		Username:      username,
		Role:          "instructor",
		Status:        "active",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	eventBody, err := json.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	instructorRepo.On("FindByInstructorID", ctx, userID).Return(instructors, nil).Once()

	handler := handlers.NewUserUpdatedHandler(instructorRepo, logger)

	err = handler.Handle(eventBody)
	require.NoError(t, err)

	instructorRepo.AssertNotCalled(t, "Update", ctx, mock.Anything)
	instructorRepo.AssertExpectations(t)
}

func TestUserUpdatedHandler_UpdateError(t *testing.T) {
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	logger := logger.NewNop()

	userID := "user-123"
	newUsername := "new_username"
	oldUsername := "old_username"

	instructor := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-1",
		CourseOfferingID:   "offering-1",
		InstructorID:       userID,
		InstructorUsername: oldUsername,
	}

	instructors := []*entities.CourseOfferingInstructor{instructor}

	event := events.UserUpdatedEvent{
		ID:            userID,
		Email:         "user@example.com",
		Username:      newUsername,
		Role:          "instructor",
		Status:        "active",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	eventBody, err := json.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	instructorRepo.On("FindByInstructorID", ctx, userID).Return(instructors, nil).Once()
	instructorRepo.On("Update", ctx, mock.Anything).Return(assert.AnError).Once()

	handler := handlers.NewUserUpdatedHandler(instructorRepo, logger)

	err = handler.Handle(eventBody)
	require.Error(t, err)

	instructorRepo.AssertExpectations(t)
}

