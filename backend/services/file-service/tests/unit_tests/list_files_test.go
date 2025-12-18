package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListFiles_Success(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	logger := logger.NewNop()

	file1 := entities.NewFile("file1.txt", "stored1.txt", "bucket", "text/plain", 100, "user-123", []string{})
	file2 := entities.NewFile("file2.txt", "stored2.txt", "bucket", "text/plain", 200, "user-123", []string{})

	limit := 10
	query := repositories.FileQuery{
		Limit: &limit,
	}

	result := &repositories.FileQueryResult{
		Files: []*entities.File{file1, file2},
		Total: 2,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()

	uc := usecases.NewListFilesUseCase(repo, logger, "http://localhost:3000")

	input := usecases.ListFilesInput{
		Limit: &limit,
	}

	files, total, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, files, 2)
	repo.AssertExpectations(t)
}

func TestListFiles_WithFilters(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	logger := logger.NewNop()

	file := entities.NewFile("file.txt", "stored.txt", "bucket", "text/plain", 100, "user-123", []string{"tag1"})

	uploadedBy := "user-123"
	limit := 10
	tags := []string{"tag1"}
	query := repositories.FileQuery{
		UploadedBy: &uploadedBy,
		Tags:       tags,
		Limit:      &limit,
	}

	result := &repositories.FileQueryResult{
		Files: []*entities.File{file},
		Total: 1,
	}

	repo.On("Find", mock.Anything, query).Return(result, nil).Once()

	uc := usecases.NewListFilesUseCase(repo, logger, "http://localhost:3000")

	input := usecases.ListFilesInput{
		UploadedBy: &uploadedBy,
		Tags:       tags,
		Limit:      &limit,
	}

	files, total, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, files, 1)
	repo.AssertExpectations(t)
}

func TestListFiles_RepositoryError(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	logger := logger.NewNop()

	limit := 10
	query := repositories.FileQuery{
		Limit: &limit,
	}

	repo.On("Find", mock.Anything, query).Return(nil, errors.New("database error")).Once()

	uc := usecases.NewListFilesUseCase(repo, logger, "http://localhost:3000")

	input := usecases.ListFilesInput{
		Limit: &limit,
	}

	_, _, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	repo.AssertExpectations(t)
}

