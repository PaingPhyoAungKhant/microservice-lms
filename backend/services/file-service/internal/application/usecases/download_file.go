package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type DownloadFileInput struct {
	FileID     string
	UserID     string
	BucketName *string
}

type DownloadFileOutput struct {
	Reader        io.ReadCloser
	ContentType   string
	ContentLength int64
	Filename      string
}

type DownloadFileUseCase struct {
	fileRepo repositories.FileRepository
	storage  *storage.MinIOClient
	logger   *logger.Logger
}

func NewDownloadFileUseCase(fileRepo repositories.FileRepository, storage *storage.MinIOClient, logger *logger.Logger) *DownloadFileUseCase {
	return &DownloadFileUseCase{
		fileRepo: fileRepo,
		storage:  storage,
		logger:   logger,
	}
}

func (uc *DownloadFileUseCase) Execute(ctx context.Context, input DownloadFileInput) (*DownloadFileOutput, error) {
	file, err := uc.fileRepo.FindByID(ctx, input.FileID)
	if err != nil {
		uc.logger.Error("file not found", zap.String("file_id", input.FileID), zap.Error(err))
		return nil, ErrFileNotFound
	}

	if file.IsDeleted() {
		return nil, ErrFileNotFound
	}

	if input.BucketName != nil && *input.BucketName != file.BucketName {
		uc.logger.Error("bucket mismatch",
			zap.String("file_id", input.FileID),
			zap.String("expected_bucket", file.BucketName),
			zap.String("provided_bucket", *input.BucketName),
		)
		return nil, ErrFileNotFound
	}

	bucketName := file.BucketName
	if input.BucketName != nil {
		bucketName = *input.BucketName
	}

	reader, err := uc.storage.DownloadFile(ctx, bucketName, file.StoredFilename)
	if err != nil {
		uc.logger.Error("failed to download file from MinIO",
			zap.String("file_id", input.FileID),
			zap.String("bucket", file.BucketName),
			zap.String("object", file.StoredFilename),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	uc.logger.Info("file downloaded",
		zap.String("file_id", input.FileID),
		zap.String("filename", file.OriginalFilename),
	)

	return &DownloadFileOutput{
		Reader:        reader,
		ContentType:   file.MimeType,
		ContentLength: file.SizeBytes,
		Filename:      file.OriginalFilename,
	}, nil
}

