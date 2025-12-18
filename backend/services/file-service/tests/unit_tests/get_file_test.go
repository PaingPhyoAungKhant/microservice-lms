package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetFile_Success(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	logger := logger.NewNop()

	file := entities.NewFile("test.txt", "stored.txt", "bucket", "text/plain", 100, "user-123", []string{})

	repo.On("FindByID", mock.Anything, "file-id").Return(file, nil).Once()

	uc := usecases.NewGetFileUseCase(repo, logger, "http://localhost:3000")

	result, err := uc.Execute(context.Background(), usecases.GetFileInput{
		FileID: "file-id",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, file.ID, result.ID)
	assert.Equal(t, file.OriginalFilename, result.OriginalFilename)
	repo.AssertExpectations(t)
}

func TestGetFile_NotFound(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	logger := logger.NewNop()

	repo.On("FindByID", mock.Anything, "file-id").Return(nil, errors.New("not found")).Once()

	uc := usecases.NewGetFileUseCase(repo, logger, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.GetFileInput{
		FileID: "file-id",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertExpectations(t)
}

