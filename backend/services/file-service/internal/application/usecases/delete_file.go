package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)


type DeleteFileInput struct {
	FileID   string
	UserID   string
}

type DeleteFileUseCase struct {
	fileRepo repositories.FileRepository
	storage  *storage.MinIOClient
	logger   *logger.Logger
}

func NewDeleteFileUseCase(fileRepo repositories.FileRepository, storage *storage.MinIOClient, logger *logger.Logger) *DeleteFileUseCase {
	return &DeleteFileUseCase{
		fileRepo: fileRepo,
		storage:  storage,
		logger:   logger,
	}
}

func (uc *DeleteFileUseCase) Execute(ctx context.Context, input DeleteFileInput) error {
	file, err := uc.fileRepo.FindByID(ctx, input.FileID)
	if err != nil {
		uc.logger.Error("file not found", zap.String("file_id", input.FileID), zap.Error(err))
		return ErrFileNotFound
	}

	if file.IsDeleted() {
		return ErrFileNotFound
	}




	if err := uc.fileRepo.SoftDelete(ctx, input.FileID); err != nil {
		return err
	}

	if err := uc.storage.DeleteFile(ctx, file.BucketName, file.StoredFilename); err != nil {
		uc.logger.Warn("failed to delete file from MinIO",
			zap.String("file_id", input.FileID),
			zap.String("bucket", file.BucketName),
			zap.String("object", file.StoredFilename),
			zap.Error(err),
		)

	}

	uc.logger.Info("file deleted",
		zap.String("file_id", input.FileID),
		zap.String("filename", file.OriginalFilename),
	)

	return nil
}

