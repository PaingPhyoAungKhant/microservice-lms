package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

var (
	ErrZoomRecordingNotFound = errors.New("zoom recording not found")
)

type CreateZoomRecordingUseCase struct {
	recordingRepo repositories.ZoomRecordingRepository
	meetingRepo   repositories.ZoomMeetingRepository
}

func NewCreateZoomRecordingUseCase(
	recordingRepo repositories.ZoomRecordingRepository,
	meetingRepo repositories.ZoomMeetingRepository,
) *CreateZoomRecordingUseCase {
	return &CreateZoomRecordingUseCase{
		recordingRepo: recordingRepo,
		meetingRepo:   meetingRepo,
	}
}

func (uc *CreateZoomRecordingUseCase) Execute(ctx context.Context, input dtos.CreateZoomRecordingInput) (*dtos.ZoomRecordingDTO, error) {
	meeting, err := uc.meetingRepo.FindByID(ctx, input.ZoomMeetingID)
	if err != nil {
		return nil, err
	}
	if meeting == nil {
		return nil, ErrZoomMeetingNotFound
	}

	recording := entities.NewZoomRecording(
		input.ZoomMeetingID,
		input.FileID,
		input.RecordingType,
		input.RecordingStartTime,
		input.RecordingEndTime,
		input.FileSize,
	)

	if err := uc.recordingRepo.Create(ctx, recording); err != nil {
		return nil, err
	}

	var dto dtos.ZoomRecordingDTO
	dto.FromEntity(recording)
	return &dto, nil
}

