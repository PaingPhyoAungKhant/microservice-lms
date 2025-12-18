package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourse_Success(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"
	thumbnailID := "thumbnail-123"
	input := dtos.CreateCourseInput{
		Name:        "Introduction to Go",
		Description: "Learn Go programming",
		ThumbnailID: &thumbnailID,
	}

	var createdCourse *entities.Course
	repo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.Course) bool {
		createdCourse = c
		return true
	})).Return(nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, mock.Anything).Return([]*entities.CourseCategory{}, nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)

	dto, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, createdCourse)
	require.Equal(t, input.Name, createdCourse.Name)
	require.Equal(t, input.Description, createdCourse.Description)
	require.Equal(t, thumbnailID, *createdCourse.ThumbnailID)

	assertDTOEqualCourse(t, dto, createdCourse)
	assert.Contains(t, dto.ThumbnailURL, apiGatewayURL)
	assert.Contains(t, dto.ThumbnailURL, "/api/v1/buckets/course-thumbnails/files/")
	assert.Contains(t, dto.ThumbnailURL, thumbnailID)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdCourse.ID, event.ID)
	assert.Equal(t, createdCourse.Name, event.Name)
	assert.Equal(t, createdCourse.Description, event.Description)
	assert.Equal(t, createdCourse.ThumbnailID, event.ThumbnailID)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourse_WithoutThumbnail(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"
	input := dtos.CreateCourseInput{
		Name:        "Introduction to Go",
		Description: "Learn Go programming",
	}

	var createdCourse *entities.Course
	repo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.Course) bool {
		createdCourse = c
		return true
	})).Return(nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, mock.Anything).Return([]*entities.CourseCategory{}, nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateCourseUseCase(repo, courseCategoryRepo, categoryRepo, publisher, logger, apiGatewayURL)

	dto, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, createdCourse)
	require.Nil(t, createdCourse.ThumbnailID)

	assertDTOEqualCourse(t, dto, createdCourse)
	assert.Empty(t, dto.ThumbnailURL)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdCourse.ID, event.ID)
	assert.Nil(t, event.ThumbnailID)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

