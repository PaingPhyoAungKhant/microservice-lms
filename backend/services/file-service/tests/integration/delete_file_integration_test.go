package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestDeleteFile_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	fileRepo := SetupFileRepository(db)
	logger := SetupTestLogger()

	storageClient := &storage.MinIOClient{}
	deleteUC := usecases.NewDeleteFileUseCase(fileRepo, storageClient, logger)

	err := deleteUC.Execute(context.Background(), usecases.DeleteFileInput{
		FileID: "00000000-0000-0000-0000-000000000000",
		UserID: uuid.New().String(),
	})

	require.Error(t, err)
	assert.Equal(t, usecases.ErrFileNotFound, err)
}

