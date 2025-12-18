package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockSectionModuleRepository struct {
	mock.Mock
}

func (m *MockSectionModuleRepository) Create(ctx context.Context, module *entities.SectionModule) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockSectionModuleRepository) FindByID(ctx context.Context, id string) (*entities.SectionModule, error) {
	args := m.Called(ctx, id)
	module, _ := args.Get(0).(*entities.SectionModule)
	return module, args.Error(1)
}

func (m *MockSectionModuleRepository) FindBySectionID(ctx context.Context, sectionID string) ([]*entities.SectionModule, error) {
	args := m.Called(ctx, sectionID)
	modules, _ := args.Get(0).([]*entities.SectionModule)
	return modules, args.Error(1)
}

func (m *MockSectionModuleRepository) Update(ctx context.Context, module *entities.SectionModule) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockSectionModuleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSectionModuleRepository) DeleteBySectionID(ctx context.Context, sectionID string) error {
	args := m.Called(ctx, sectionID)
	return args.Error(0)
}

