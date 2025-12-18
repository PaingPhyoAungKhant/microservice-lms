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

func TestUpdateCourseSection_Success(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	existingSection := &entities.CourseSection{
		ID:               sectionID,
		CourseOfferingID: "offering-123",
		Name:             "Introduction",
		Description:      "Old description",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	input := dtos.UpdateCourseSectionInput{
		Name:        "Advanced Topics",
		Description: "Updated description",
		Order:       2,
	}

	sectionRepo.On("FindByID", mock.Anything, sectionID).Return(existingSection, nil).Once()
	sectionRepo.On("Update", mock.Anything, mock.MatchedBy(func(s *entities.CourseSection) bool {
		return s.ID == sectionID && s.Name == input.Name && s.Description == input.Description
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseSectionUpdated, mock.Anything).Return(nil).Once()

	uc := usecases.NewUpdateCourseSectionUseCase(sectionRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), sectionID, input)
	require.NoError(t, err)
	require.Equal(t, input.Name, existingSection.Name)
	require.Equal(t, input.Description, existingSection.Description)
	require.Equal(t, input.Order, existingSection.Order)

	assertDTOEqualCourseSection(t, dto, existingSection)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseSectionUpdatedEvent)
	require.True(t, ok)
	assert.Equal(t, existingSection.ID, event.ID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, input.Order, event.Order)

	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourseSection_SectionNotFound(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	input := dtos.UpdateCourseSectionInput{
		Name:        "Advanced Topics",
		Description: "Updated description",
		Order:       2,
	}

	sectionRepo.On("FindByID", mock.Anything, sectionID).Return((*entities.CourseSection)(nil), nil).Once()

	uc := usecases.NewUpdateCourseSectionUseCase(sectionRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), sectionID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseSectionNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	sectionRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourseSection_RepositoryError(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	existingSection := &entities.CourseSection{
		ID:               sectionID,
		CourseOfferingID: "offering-123",
		Name:             "Introduction",
		Description:      "Old description",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	input := dtos.UpdateCourseSectionInput{
		Name:        "Advanced Topics",
		Description: "Updated description",
		Order:       2,
	}

	sectionRepo.On("FindByID", mock.Anything, sectionID).Return(existingSection, nil).Once()
	sectionRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewUpdateCourseSectionUseCase(sectionRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), sectionID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

