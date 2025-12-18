package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *entities.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) FindByID(ctx context.Context, id string) (*entities.File, error) {
	args := m.Called(ctx, id)
	file, _ := args.Get(0).(*entities.File)
	return file, args.Error(1)
}

func (m *MockFileRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.File, error) {
	args := m.Called(ctx, userID, limit, offset)
	files, _ := args.Get(0).([]*entities.File)
	return files, args.Error(1)
}

func (m *MockFileRepository) FindByTags(ctx context.Context, tags []string, limit, offset int) ([]*entities.File, error) {
	args := m.Called(ctx, tags, limit, offset)
	files, _ := args.Get(0).([]*entities.File)
	return files, args.Error(1)
}

func (m *MockFileRepository) Find(ctx context.Context, query repositories.FileQuery) (*repositories.FileQueryResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repositories.FileQueryResult)
	return result, args.Error(1)
}

func (m *MockFileRepository) Update(ctx context.Context, file *entities.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

