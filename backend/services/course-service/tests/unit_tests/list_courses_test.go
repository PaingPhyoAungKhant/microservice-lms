package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCourses_Success(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	course1 := entities.NewCourse("Course 1", "Description 1", nil)
	course2 := entities.NewCourse("Course 2", "Description 2", nil)

	limit := 10
	query := repositories.CourseQuery{
		Limit: &limit,
	}

	result := &repositories.CourseQueryResult{
		Courses: []*entities.Course{course1, course2},
		Total:   2,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course1.ID).Return([]*entities.CourseCategory{}, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course2.ID).Return([]*entities.CourseCategory{}, nil).Once()

	uc := usecases.NewListCoursesUseCase(repo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	input := usecases.ListCoursesInput{
		Limit: &limit,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.Total)
	assert.Len(t, output.Courses, 2)

	repo.AssertExpectations(t)
	courseCategoryRepo.AssertExpectations(t)
}

func TestListCourses_WithCategoryFilter(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	categoryID := "category-123"
	limit := 10
	query := repositories.CourseQuery{
		CategoryID: &categoryID,
		Limit:      &limit,
	}

	course := entities.NewCourse("Course 1", "Description 1", nil)
	result := &repositories.CourseQueryResult{
		Courses: []*entities.Course{course},
		Total:   1,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course.ID).Return([]*entities.CourseCategory{}, nil).Once()

	uc := usecases.NewListCoursesUseCase(repo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	input := usecases.ListCoursesInput{
		CategoryID: &categoryID,
		Limit:      &limit,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 1, output.Total)
	assert.Len(t, output.Courses, 1)

	repo.AssertExpectations(t)
	courseCategoryRepo.AssertExpectations(t)
}

func TestListCourses_WithSearch(t *testing.T) {
	repo := new(mocks.MockCourseRepository)
	courseCategoryRepo := new(mocks.MockCourseCategoryRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()
	apiGatewayURL := "http://localhost:3000"

	searchQuery := "Go"
	limit := 10
	query := repositories.CourseQuery{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	course := entities.NewCourse("Introduction to Go", "Learn Go programming", nil)
	result := &repositories.CourseQueryResult{
		Courses: []*entities.Course{course},
		Total:   1,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()
	courseCategoryRepo.On("FindByCourseID", mock.Anything, course.ID).Return([]*entities.CourseCategory{}, nil).Once()

	uc := usecases.NewListCoursesUseCase(repo, courseCategoryRepo, categoryRepo, logger, apiGatewayURL)

	input := usecases.ListCoursesInput{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 1, output.Total)
	assert.Len(t, output.Courses, 1)

	repo.AssertExpectations(t)
	courseCategoryRepo.AssertExpectations(t)
}

