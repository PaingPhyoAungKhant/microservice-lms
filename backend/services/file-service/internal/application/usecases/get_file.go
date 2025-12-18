package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)



type GetFileInput struct {
	FileID string
}

type GetFileUseCase struct {
	fileRepo repositories.FileRepository
	logger   *logger.Logger
	apiGatewayURL string
}

func NewGetFileUseCase(fileRepo repositories.FileRepository, logger *logger.Logger, apiGatewayURL string) *GetFileUseCase {
	return &GetFileUseCase{
		fileRepo: fileRepo,
		logger:   logger,
		apiGatewayURL: apiGatewayURL,
	}
}

func (uc *GetFileUseCase) Execute(ctx context.Context, input GetFileInput) (*dtos.FileDTO, error) {
	file, err := uc.fileRepo.FindByID(ctx, input.FileID)
	if err != nil {
		uc.logger.Error("file not found", zap.String("file_id", input.FileID), zap.Error(err))
		return nil, ErrFileNotFound
	}

	var dto dtos.FileDTO
	dto.FromEntity(file, uc.apiGatewayURL)
	return &dto, nil
}

