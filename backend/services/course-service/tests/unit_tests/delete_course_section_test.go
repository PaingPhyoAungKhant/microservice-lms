package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestDeleteCourseSection_FindByIDError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "section-1").Return((*entities.CourseSection)(nil), expectedErr).Once()

	uc := usecases.NewDeleteCourseSectionUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteCourseSectionInput{SectionID: "section-1"})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteCourseSection_SectionNotFound(t *testing.T) {
	repo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "section-1").Return((*entities.CourseSection)(nil), nil).Once()

	uc := usecases.NewDeleteCourseSectionUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteCourseSectionInput{SectionID: "section-1"})
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseSectionNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteCourseSection_Success(t *testing.T) {
	section := &entities.CourseSection{
		ID:               "section-1",
		CourseOfferingID: "offering-1",
		Name:             "Test Section",
		Description:      "Test Description",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	repo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)

	repo.On("FindByID", mock.Anything, section.ID).Return(section, nil).Once()
	repo.On("Delete", mock.Anything, section.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseSectionDeleted, mock.Anything).Return(nil).Once()

	uc := usecases.NewDeleteCourseSectionUseCase(repo, publisher, logger.NewNop())

	out, err := uc.Execute(context.Background(), usecases.DeleteCourseSectionInput{SectionID: section.ID})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotEmpty(t, out.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.CourseSectionDeletedEvent)
	require.True(t, ok)
	assert.Equal(t, section.ID, event.ID)
	assert.WithinDuration(t, time.Now(), event.DeletedAt, 2*time.Second)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

