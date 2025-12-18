package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestUpdateCourse_RepoError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockCourseRepository)
	repo.On("FindByID", mock.Anything, "course-123").Return((*entities.Course)(nil), expectedErr).Once()
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	uc := usecases.NewUpdateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)
	_, err := uc.Execute(context.Background(), usecases.UpdateCourseInput{CourseID: "course-123"}, dtos.UpdateCourseInput{})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourse_CourseNotFound(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	repo.On("FindByID", mock.Anything, "course-123").Return((*entities.Course)(nil), nil).Once()
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	uc := usecases.NewUpdateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)
	_, err := uc.Execute(context.Background(), usecases.UpdateCourseInput{CourseID: "course-123"}, dtos.UpdateCourseInput{})
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourse_RepositoryError(t *testing.T) {
	now := time.Now().Add(-time.Hour)
	course := &entities.Course{
		ID:          "course-123",
		Name:        "Old Course",
		Description: "Old Description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	repo.On("FindByID", mock.Anything, course.ID).Return(course, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course.ID).Return([]*entities.CourseCategory{}, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewUpdateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)
	input := dtos.UpdateCourseInput{
		Name:        "Updated Course",
		Description: "Updated Description",
	}
	_, err := uc.Execute(context.Background(), usecases.UpdateCourseInput{CourseID: course.ID}, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateCourse_Success(t *testing.T) {
	now := time.Now().Add(-time.Hour)
	thumbnailID := "thumbnail-123"
	course := &entities.Course{
		ID:          "course-123",
		Name:        "Old Course",
		Description: "Old Description",
		ThumbnailID: &thumbnailID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	var updatedCourse *entities.Course
	repo.On("FindByID", mock.Anything, course.ID).Return(course, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course.ID).Return([]*entities.CourseCategory{}, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(c *entities.Course) bool {
		updatedCourse = c
		return true
	})).Return(nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course.ID).Return([]*entities.CourseCategory{}, nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseUpdated, mock.Anything).Return(nil).Once()

	newName := "Updated Course"
	newDescription := "Updated Description"
	newThumbnailID := "thumbnail-456"
	uc := usecases.NewUpdateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)
	input := dtos.UpdateCourseInput{
		Name:        newName,
		Description: newDescription,
		ThumbnailID:  &newThumbnailID,
	}
	dto, err := uc.Execute(context.Background(), usecases.UpdateCourseInput{CourseID: course.ID}, input)
	require.NoError(t, err)
	require.NotNil(t, updatedCourse)
	assert.Equal(t, newName, updatedCourse.Name)
	assert.Equal(t, newDescription, updatedCourse.Description)
	assert.Equal(t, newThumbnailID, *updatedCourse.ThumbnailID)
	assert.True(t, updatedCourse.UpdatedAt.After(now))

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.CourseUpdatedEvent)
	require.True(t, ok)
	assert.Equal(t, updatedCourse.ID, event.ID)
	assert.Equal(t, updatedCourse.Name, event.Name)
	assert.Equal(t, updatedCourse.Description, event.Description)
	assert.Equal(t, updatedCourse.ThumbnailID, event.ThumbnailID)

	assertDTOEqualCourse(t, dto, updatedCourse)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}

