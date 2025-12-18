package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

var (
	ErrFileTooLarge      = errors.New("file size exceeds maximum allowed size")
	ErrInvalidMimeType   = errors.New("invalid mime type")
	ErrUploadFailed      = errors.New("failed to upload file")
	ErrBucketRequired    = errors.New("bucket name is required")
	ErrFileRequired      = errors.New("file is required")
)

const (
	MaxFileSize = 1024 * 1024 * 1024 
)

type UploadFileInput struct {
	File        io.Reader
	Filename    string
	MimeType    string
	Size        int64
	UploadedBy  string
	BucketName  string
	Tags        []string
}

type UploadFileUseCase struct {
	fileRepo repositories.FileRepository
	storage  *storage.MinIOClient
	logger   *logger.Logger
	apiGatewayURL string
}

func NewUploadFileUseCase(fileRepo repositories.FileRepository, storage *storage.MinIOClient, logger *logger.Logger, apiGatewayURL string) *UploadFileUseCase {
	return &UploadFileUseCase{
		fileRepo: fileRepo,
		storage:  storage,
		logger:   logger,
		apiGatewayURL: apiGatewayURL,
	}
}

func (uc *UploadFileUseCase) Execute(ctx context.Context, input UploadFileInput) (*dtos.FileDTO, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	bucketName := input.BucketName
	if bucketName == "" {
		bucketName = uc.determineBucket(input.MimeType, input.Tags)
	}

	ext := filepath.Ext(input.Filename)
	storedFilename := uuid.NewString() + ext

	if err := uc.storage.UploadFile(ctx, bucketName, storedFilename, input.File, input.Size, input.MimeType); err != nil {
		uc.logger.Error("failed to upload file to MinIO", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	file := entities.NewFile(
		input.Filename,
		storedFilename,
		bucketName,
		input.MimeType,
		input.Size,
		input.UploadedBy,
		input.Tags,
	)

	if err := uc.fileRepo.Create(ctx, file); err != nil {
		uc.storage.DeleteFile(ctx, bucketName, storedFilename)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	uc.logger.Info("file uploaded successfully",
		zap.String("file_id", file.ID),
		zap.String("filename", input.Filename),
		zap.String("bucket", bucketName),
	)

	var dto dtos.FileDTO
	dto.FromEntity(file, uc.apiGatewayURL)
	return &dto, nil
}

func (uc *UploadFileUseCase) validateInput(input UploadFileInput) error {
	if input.File == nil {
		return ErrFileRequired
	}

	if input.Size > MaxFileSize {
		return fmt.Errorf("%w: maximum size is %d bytes", ErrFileTooLarge, MaxFileSize)
	}

	if input.Filename == "" {
		return errors.New("filename is required")
	}

	if input.MimeType == "" {
		return errors.New("mime type is required")
	}

	if input.UploadedBy == "" {
		return errors.New("uploaded_by is required")
	}

	return nil
}

func (uc *UploadFileUseCase) determineBucket(mimeType string, tags []string) string {
	for _, tag := range tags {
		switch tag {
		case "thumbnail":
			return "course-thumbnails"
		case "video":
			return "course-videos"
		case "recording":
			return "zoom-recordings"
		}
	}

	if strings.HasPrefix(mimeType, "image/") {
		return "course-thumbnails"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "course-videos"
	}

	return "general-files"
}

