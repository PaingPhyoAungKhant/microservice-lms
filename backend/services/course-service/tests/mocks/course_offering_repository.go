package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockCourseOfferingRepository struct {
	mock.Mock
}

func (m *MockCourseOfferingRepository) Create(ctx context.Context, offering *entities.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}

func (m *MockCourseOfferingRepository) FindByID(ctx context.Context, id string) (*entities.CourseOffering, error) {
	args := m.Called(ctx, id)
	offering, _ := args.Get(0).(*entities.CourseOffering)
	return offering, args.Error(1)
}

func (m *MockCourseOfferingRepository) FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseOffering, error) {
	args := m.Called(ctx, courseID)
	offerings, _ := args.Get(0).([]*entities.CourseOffering)
	return offerings, args.Error(1)
}

func (m *MockCourseOfferingRepository) Find(ctx context.Context, query repositories.CourseOfferingQuery) (*repositories.CourseOfferingQueryResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repositories.CourseOfferingQueryResult)
	return result, args.Error(1)
}

func (m *MockCourseOfferingRepository) Update(ctx context.Context, offering *entities.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}

func (m *MockCourseOfferingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

