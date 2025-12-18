package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockCourseCategoryRepository struct {
	mock.Mock
}

func (m *MockCourseCategoryRepository) Create(ctx context.Context, courseCategory *entities.CourseCategory) error {
	args := m.Called(ctx, courseCategory)
	return args.Error(0)
}

func (m *MockCourseCategoryRepository) FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseCategory, error) {
	args := m.Called(ctx, courseID)
	courseCategories, _ := args.Get(0).([]*entities.CourseCategory)
	return courseCategories, args.Error(1)
}

func (m *MockCourseCategoryRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*entities.CourseCategory, error) {
	args := m.Called(ctx, categoryID)
	courseCategories, _ := args.Get(0).([]*entities.CourseCategory)
	return courseCategories, args.Error(1)
}

func (m *MockCourseCategoryRepository) Delete(ctx context.Context, courseID, categoryID string) error {
	args := m.Called(ctx, courseID, categoryID)
	return args.Error(0)
}

func (m *MockCourseCategoryRepository) DeleteByCourseID(ctx context.Context, courseID string) error {
	args := m.Called(ctx, courseID)
	return args.Error(0)
}

