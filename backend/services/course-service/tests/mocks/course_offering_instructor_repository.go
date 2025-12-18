package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockCourseOfferingInstructorRepository struct {
	mock.Mock
}

func (m *MockCourseOfferingInstructorRepository) Create(ctx context.Context, instructor *entities.CourseOfferingInstructor) error {
	args := m.Called(ctx, instructor)
	return args.Error(0)
}

func (m *MockCourseOfferingInstructorRepository) FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseOfferingInstructor, error) {
	args := m.Called(ctx, offeringID)
	instructors, _ := args.Get(0).([]*entities.CourseOfferingInstructor)
	return instructors, args.Error(1)
}

func (m *MockCourseOfferingInstructorRepository) FindByInstructorID(ctx context.Context, instructorID string) ([]*entities.CourseOfferingInstructor, error) {
	args := m.Called(ctx, instructorID)
	instructors, _ := args.Get(0).([]*entities.CourseOfferingInstructor)
	return instructors, args.Error(1)
}

func (m *MockCourseOfferingInstructorRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseOfferingInstructorRepository) DeleteByOfferingID(ctx context.Context, offeringID string) error {
	args := m.Called(ctx, offeringID)
	return args.Error(0)
}

func (m *MockCourseOfferingInstructorRepository) Update(ctx context.Context, instructor *entities.CourseOfferingInstructor) error {
	args := m.Called(ctx, instructor)
	return args.Error(0)
}

