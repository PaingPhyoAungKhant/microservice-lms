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

func TestDeleteCourse_FindByIDError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "course-1").Return((*entities.Course)(nil), expectedErr).Once()

	uc := usecases.NewDeleteCourseUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteCourseInput{CourseID: "course-1"})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteCourse_CourseNotFound(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "course-1").Return((*entities.Course)(nil), nil).Once()

	uc := usecases.NewDeleteCourseUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteCourseInput{CourseID: "course-1"})
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteCourse_Success(t *testing.T) {
	course := &entities.Course{
		ID:          "course-1",
		Name:        "Test Course",
		Description: "Test Description",
	}

	repo := new(mocks.MockCourseRepository)
	publisher := new(sharedMocks.MockPublisher)

	repo.On("FindByID", mock.Anything, course.ID).Return(course, nil).Once()
	repo.On("Delete", mock.Anything, course.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseDeleted, mock.Anything).Return(nil).Once()

	uc := usecases.NewDeleteCourseUseCase(repo, publisher, logger.NewNop())

	out, err := uc.Execute(context.Background(), usecases.DeleteCourseInput{CourseID: course.ID})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotEmpty(t, out.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.CourseDeletedEvent)
	require.True(t, ok)
	assert.Equal(t, course.ID, event.ID)
	assert.WithinDuration(t, time.Now(), event.DeletedAt, 2*time.Second)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

