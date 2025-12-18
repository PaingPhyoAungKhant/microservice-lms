package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) Create(ctx context.Context, course *entities.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) FindByID(ctx context.Context, id string) (*entities.Course, error) {
	args := m.Called(ctx, id)
	course, _ := args.Get(0).(*entities.Course)
	return course, args.Error(1)
}

func (m *MockCourseRepository) Find(ctx context.Context, query repositories.CourseQuery) (*repositories.CourseQueryResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repositories.CourseQueryResult)
	return result, args.Error(1)
}

func (m *MockCourseRepository) Update(ctx context.Context, course *entities.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

