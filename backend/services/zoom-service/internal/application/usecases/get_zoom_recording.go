package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type GetZoomRecordingUseCase struct {
	recordingRepo repositories.ZoomRecordingRepository
}

func NewGetZoomRecordingUseCase(recordingRepo repositories.ZoomRecordingRepository) *GetZoomRecordingUseCase {
	return &GetZoomRecordingUseCase{
		recordingRepo: recordingRepo,
	}
}

func (uc *GetZoomRecordingUseCase) Execute(ctx context.Context, id string) (*dtos.ZoomRecordingDTO, error) {
	recording, err := uc.recordingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if recording == nil {
		return nil, ErrZoomRecordingNotFound
	}

	var dto dtos.ZoomRecordingDTO
	dto.FromEntity(recording)
	return &dto, nil
}

