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

func TestListFiles_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	fileRepo := SetupFileRepository(db)
	logger := SetupTestLogger()

	ctx := context.Background()

	userID := uuid.New().String()

	file1 := entities.NewFile(
		"file1.txt",
		"stored-"+uuid.New().String()+".txt",
		"general-files",
		"text/plain",
		100,
		userID,
		[]string{},
	)
	file2 := entities.NewFile(
		"file2.txt",
		"stored-"+uuid.New().String()+".txt",
		"general-files",
		"text/plain",
		200,
		userID,
		[]string{},
	)

	err := fileRepo.Create(ctx, file1)
	require.NoError(t, err)
	err = fileRepo.Create(ctx, file2)
	require.NoError(t, err)

	listUC := usecases.NewListFilesUseCase(fileRepo, logger, "http://localhost:3000")

	limit := 10
	result, total, err := listUC.Execute(ctx, usecases.ListFilesInput{
		Limit: &limit,
	})

	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 2)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestListFiles_Integration_WithFilters(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	fileRepo := SetupFileRepository(db)
	logger := SetupTestLogger()

	ctx := context.Background()

	userID := uuid.New().String()
	tags := []string{"thumbnail", "course"}

	file := entities.NewFile(
		"image.jpg",
		"stored-"+uuid.New().String()+".jpg",
		"course-thumbnails",
		"image/jpeg",
		5000,
		userID,
		tags,
	)

	err := fileRepo.Create(ctx, file)
	require.NoError(t, err)

	listUC := usecases.NewListFilesUseCase(fileRepo, logger, "http://localhost:3000")

	limit := 10
	result, total, err := listUC.Execute(ctx, usecases.ListFilesInput{
		UploadedBy: &userID,
		Tags:       tags,
		Limit:      &limit,
	})

	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 1)
	assert.GreaterOrEqual(t, len(result), 1)
}

