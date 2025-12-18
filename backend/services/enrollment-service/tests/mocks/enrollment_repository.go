package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) Create(ctx context.Context, enrollment *entities.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) FindByID(ctx context.Context, id string) (*entities.Enrollment, error) {
	args := m.Called(ctx, id)
	enrollment, _ := args.Get(0).(*entities.Enrollment)
	return enrollment, args.Error(1)
}

func (m *MockEnrollmentRepository) Update(ctx context.Context, enrollment *entities.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Find(ctx context.Context, query repositories.EnrollmentQuery) (*repositories.EnrollmentQueryResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repositories.EnrollmentQueryResult)
	return result, args.Error(1)
}

func (m *MockEnrollmentRepository) UpdateStudentUsername(ctx context.Context, studentID, username string) error {
	args := m.Called(ctx, studentID, username)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) UpdateCourseName(ctx context.Context, courseID, courseName string) error {
	args := m.Called(ctx, courseID, courseName)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) UpdateCourseOfferingName(ctx context.Context, courseOfferingID, courseOfferingName string) error {
	args := m.Called(ctx, courseOfferingID, courseOfferingName)
	return args.Error(0)
}

