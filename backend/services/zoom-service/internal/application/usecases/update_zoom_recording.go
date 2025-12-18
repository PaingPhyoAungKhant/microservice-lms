package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type UpdateZoomRecordingUseCase struct {
	recordingRepo repositories.ZoomRecordingRepository
}

func NewUpdateZoomRecordingUseCase(recordingRepo repositories.ZoomRecordingRepository) *UpdateZoomRecordingUseCase {
	return &UpdateZoomRecordingUseCase{
		recordingRepo: recordingRepo,
	}
}

func (uc *UpdateZoomRecordingUseCase) Execute(ctx context.Context, id string, input dtos.UpdateZoomRecordingInput) (*dtos.ZoomRecordingDTO, error) {
	recording, err := uc.recordingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if recording == nil {
		return nil, ErrZoomRecordingNotFound
	}

	recording.Update(input.RecordingType, input.RecordingStartTime, input.RecordingEndTime, input.FileSize)

	if err := uc.recordingRepo.Update(ctx, recording); err != nil {
		return nil, err
	}

	var dto dtos.ZoomRecordingDTO
	dto.FromEntity(recording)
	return &dto, nil
}

