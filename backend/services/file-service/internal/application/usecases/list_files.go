package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type ListFilesInput struct {
	UploadedBy *string
	Tags       []string
	MimeType   *string
	BucketName *string
	Limit      *int
	Offset     *int
	SortColumn *string
	SortDirection *string
}

type ListFilesUseCase struct {
	fileRepo repositories.FileRepository
	logger   *logger.Logger
	apiGatewayURL string
}

func NewListFilesUseCase(fileRepo repositories.FileRepository, logger *logger.Logger, apiGatewayURL string) *ListFilesUseCase {
	return &ListFilesUseCase{
		fileRepo: fileRepo,
		logger:   logger,
		apiGatewayURL: apiGatewayURL,
	}
}

func (uc *ListFilesUseCase) Execute(ctx context.Context, input ListFilesInput) ([]*dtos.FileDTO, int, error) {
	query := repositories.FileQuery{
		UploadedBy: input.UploadedBy,
		Tags:       input.Tags,
		MimeType:   input.MimeType,
		BucketName: input.BucketName,
		Limit:      input.Limit,
		Offset:     input.Offset,
		SortColumn: input.SortColumn,
		SortDirection: input.SortDirection,
	}

	result, err := uc.fileRepo.Find(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	fileDTOs := make([]*dtos.FileDTO, 0, len(result.Files))
	for _, file := range result.Files {
		var dto dtos.FileDTO
		dto.FromEntity(file, uc.apiGatewayURL)
		fileDTOs = append(fileDTOs, &dto)
	}

	return fileDTOs, result.Total, nil
}

