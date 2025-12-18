package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockCourseSectionRepository struct {
	mock.Mock
}

func (m *MockCourseSectionRepository) Create(ctx context.Context, section *entities.CourseSection) error {
	args := m.Called(ctx, section)
	return args.Error(0)
}

func (m *MockCourseSectionRepository) FindByID(ctx context.Context, id string) (*entities.CourseSection, error) {
	args := m.Called(ctx, id)
	section, _ := args.Get(0).(*entities.CourseSection)
	return section, args.Error(1)
}

func (m *MockCourseSectionRepository) FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseSection, error) {
	args := m.Called(ctx, offeringID)
	sections, _ := args.Get(0).([]*entities.CourseSection)
	return sections, args.Error(1)
}

func (m *MockCourseSectionRepository) Update(ctx context.Context, section *entities.CourseSection) error {
	args := m.Called(ctx, section)
	return args.Error(0)
}

func (m *MockCourseSectionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseSectionRepository) DeleteByOfferingID(ctx context.Context, offeringID string) error {
	args := m.Called(ctx, offeringID)
	return args.Error(0)
}

