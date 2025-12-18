package unit_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUploadFile_MissingFile(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	uc := usecases.NewUploadFileUseCase(repo, storageClient, logger, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.UploadFileInput{
		Filename:   "test.txt",
		MimeType:   "text/plain",
		Size:       100,
		UploadedBy: "user-123",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileRequired, err)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadFile_FileTooLarge(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	uc := usecases.NewUploadFileUseCase(repo, storageClient, logger, "http://localhost:3000")

	fileContent := bytes.NewReader(make([]byte, 2*1024*1024*1024))

	_, err := uc.Execute(context.Background(), usecases.UploadFileInput{
		File:       fileContent,
		Filename:   "large.txt",
		MimeType:   "text/plain",
		Size:       2 * 1024 * 1024 * 1024,
		UploadedBy: "user-123",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "file size exceeds maximum allowed size")
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadFile_MissingFilename(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	uc := usecases.NewUploadFileUseCase(repo, storageClient, logger, "http://localhost:3000")

	fileContent := bytes.NewReader([]byte("test content"))

	_, err := uc.Execute(context.Background(), usecases.UploadFileInput{
		File:       fileContent,
		MimeType:   "text/plain",
		Size:       100,
		UploadedBy: "user-123",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "filename is required")
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadFile_MissingMimeType(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	uc := usecases.NewUploadFileUseCase(repo, storageClient, logger, "http://localhost:3000")

	fileContent := bytes.NewReader([]byte("test content"))

	_, err := uc.Execute(context.Background(), usecases.UploadFileInput{
		File:       fileContent,
		Filename:   "test.txt",
		Size:       100,
		UploadedBy: "user-123",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "mime type is required")
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadFile_MissingUploadedBy(t *testing.T) {
	repo := new(mocks.MockFileRepository)
	storageClient := &storage.MinIOClient{}
	logger := logger.NewNop()

	uc := usecases.NewUploadFileUseCase(repo, storageClient, logger, "http://localhost:3000")

	fileContent := bytes.NewReader([]byte("test content"))

	_, err := uc.Execute(context.Background(), usecases.UploadFileInput{
		File:     fileContent,
		Filename: "test.txt",
		MimeType: "text/plain",
		Size:     100,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "uploaded_by is required")
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

