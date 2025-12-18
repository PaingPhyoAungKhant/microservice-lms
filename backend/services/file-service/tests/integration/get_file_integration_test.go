package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFile_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	fileRepo := SetupFileRepository(db)
	logger := SetupTestLogger()

	ctx := context.Background()

	file := entities.NewFile(
		"test.txt",
		"stored-"+uuid.New().String()+".txt",
		"general-files",
		"text/plain",
		100,
		uuid.New().String(),
		[]string{},
	)

	err := fileRepo.Create(ctx, file)
	require.NoError(t, err)

	getUC := usecases.NewGetFileUseCase(fileRepo, logger, "http://localhost:3000")

	result, err := getUC.Execute(ctx, usecases.GetFileInput{
		FileID: file.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, file.ID, result.ID)
	assert.Equal(t, file.OriginalFilename, result.OriginalFilename)
	assert.Equal(t, file.BucketName, result.BucketName)
}

func TestGetFile_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	fileRepo := SetupFileRepository(db)
	logger := SetupTestLogger()

	getUC := usecases.NewGetFileUseCase(fileRepo, logger, "http://localhost:3000")

	_, err := getUC.Execute(context.Background(), usecases.GetFileInput{
		FileID: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
}

