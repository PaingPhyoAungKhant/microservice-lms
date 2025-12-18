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

func TestUpdateCourseOffering_Success(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	existingOffering := &entities.CourseOffering{
		ID:             offeringID,
		CourseID:       "course-123",
		Name:           "Spring 2024",
		Description:    "Old description",
		OfferingType:   entities.OfferingTypeOnline,
		Status:         entities.OfferingStatusPending,
		EnrollmentCost: 50.0,
	}

	duration := "6 months"
	classTime := "Tue, Thu 2:00 PM"
	input := dtos.UpdateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Updated description",
		OfferingType:   "oncampus",
		Duration:       &duration,
		ClassTime:      &classTime,
		EnrollmentCost: 149.99,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(existingOffering, nil).Once()
	offeringRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *entities.CourseOffering) bool {
		return o.ID == offeringID && o.Name == input.Name && o.Description == input.Description
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseOfferingUpdated, mock.Anything).Return(nil).Once()

	uc := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), offeringID, input)
	require.NoError(t, err)
	require.Equal(t, input.Name, existingOffering.Name)
	require.Equal(t, input.Description, existingOffering.Description)
	require.Equal(t, entities.OfferingTypeOnCampus, existingOffering.OfferingType)
	require.Equal(t, duration, *existingOffering.Duration)
	require.Equal(t, classTime, *existingOffering.ClassTime)
	require.Equal(t, input.EnrollmentCost, existingOffering.EnrollmentCost)

	assertDTOEqualCourseOffering(t, dto, existingOffering)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseOfferingUpdatedEvent)
	require.True(t, ok)
	assert.Equal(t, existingOffering.ID, event.ID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, string(existingOffering.OfferingType), event.OfferingType)
	assert.Equal(t, input.EnrollmentCost, event.EnrollmentCost)

	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourseOffering_OfferingNotFound(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	input := dtos.UpdateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Updated description",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return((*entities.CourseOffering)(nil), nil).Once()

	uc := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseOfferingNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	offeringRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourseOffering_RepositoryError(t *testing.T) {
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	existingOffering := &entities.CourseOffering{
		ID:             offeringID,
		CourseID:       "course-123",
		Name:           "Spring 2024",
		Description:    "Old description",
		OfferingType:   entities.OfferingTypeOnline,
		Status:         entities.OfferingStatusPending,
		EnrollmentCost: 50.0,
	}

	input := dtos.UpdateCourseOfferingInput{
		Name:           "Fall 2024",
		Description:    "Updated description",
		OfferingType:   "online",
		EnrollmentCost: 99.99,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(existingOffering, nil).Once()
	offeringRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

