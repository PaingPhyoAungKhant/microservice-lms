package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type DeleteZoomRecordingUseCase struct {
	recordingRepo repositories.ZoomRecordingRepository
}

func NewDeleteZoomRecordingUseCase(recordingRepo repositories.ZoomRecordingRepository) *DeleteZoomRecordingUseCase {
	return &DeleteZoomRecordingUseCase{
		recordingRepo: recordingRepo,
	}
}

func (uc *DeleteZoomRecordingUseCase) Execute(ctx context.Context, id string) error {
	recording, err := uc.recordingRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if recording == nil {
		return ErrZoomRecordingNotFound
	}

	if err := uc.recordingRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

