package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDownloadFile_NotFound(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	repo.On("FindByID", mock.Anything, "file-id").Return(nil, errors.New("not found")).Once()

	uc := usecases.NewDownloadFileUseCase(repo, storageClient, logger)

	_, err := uc.Execute(context.Background(), usecases.DownloadFileInput{
		FileID: "file-id",
		UserID: "user-123",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertExpectations(t)
}

func TestDownloadFile_BucketMismatch(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	file := entities.NewFile("test.txt", "stored.txt", "bucket1", "text/plain", 100, "user-123", []string{})

	repo.On("FindByID", mock.Anything, "file-id").Return(file, nil).Once()

	uc := usecases.NewDownloadFileUseCase(repo, storageClient, logger)

	bucket := "bucket2"
	_, err := uc.Execute(context.Background(), usecases.DownloadFileInput{
		FileID:     "file-id",
		UserID:     "user-123",
		BucketName: &bucket,
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertExpectations(t)
}

func TestDownloadFile_FileAlreadyDeleted(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	file := entities.NewFile("test.txt", "stored.txt", "bucket", "text/plain", 100, "user-123", []string{})
	file.SoftDelete()

	repo.On("FindByID", mock.Anything, "file-id").Return(file, nil).Once()

	uc := usecases.NewDownloadFileUseCase(repo, storageClient, logger)

	_, err := uc.Execute(context.Background(), usecases.DownloadFileInput{
		FileID: "file-id",
		UserID: "user-123",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertExpectations(t)
}

