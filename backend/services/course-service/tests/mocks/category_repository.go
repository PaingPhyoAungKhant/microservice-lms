package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(ctx context.Context, id string) (*entities.Category, error) {
	args := m.Called(ctx, id)
	category, _ := args.Get(0).(*entities.Category)
	return category, args.Error(1)
}

func (m *MockCategoryRepository) FindByName(ctx context.Context, name string) (*entities.Category, error) {
	args := m.Called(ctx, name)
	category, _ := args.Get(0).(*entities.Category)
	return category, args.Error(1)
}

func (m *MockCategoryRepository) Find(ctx context.Context, query repositories.CategoryQuery) (*repositories.CategoryQueryResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repositories.CategoryQueryResult)
	return result, args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

