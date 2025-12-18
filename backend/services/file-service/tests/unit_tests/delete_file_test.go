package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)


func TestDeleteFile_NotFound(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	repo.On("FindByID", mock.Anything, "file-id").Return(nil, errors.New("not found")).Once()

	uc := usecases.NewDeleteFileUseCase(repo, storageClient, logger)

	err := uc.Execute(context.Background(), usecases.DeleteFileInput{
		FileID: "file-id",
		UserID: "user-123",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertNotCalled(t, "SoftDelete", mock.Anything, mock.Anything)
}

func TestDeleteFile_AlreadyDeleted(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	file := entities.NewFile("test.txt", "stored.txt", "bucket", "text/plain", 100, "user-123", []string{})
	now := time.Now().UTC()
	file.DeletedAt = &now

	repo.On("FindByID", mock.Anything, "file-id").Return(file, nil).Once()

	uc := usecases.NewDeleteFileUseCase(repo, storageClient, logger)

	err := uc.Execute(context.Background(), usecases.DeleteFileInput{
		FileID: "file-id",
		UserID: "user-123",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
	repo.AssertNotCalled(t, "SoftDelete", mock.Anything, mock.Anything)
}


