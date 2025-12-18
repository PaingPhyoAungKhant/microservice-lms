package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type ListZoomRecordingsUseCase struct {
	recordingRepo repositories.ZoomRecordingRepository
	meetingRepo   repositories.ZoomMeetingRepository
}

func NewListZoomRecordingsUseCase(
	recordingRepo repositories.ZoomRecordingRepository,
	meetingRepo repositories.ZoomMeetingRepository,
) *ListZoomRecordingsUseCase {
	return &ListZoomRecordingsUseCase{
		recordingRepo: recordingRepo,
		meetingRepo:   meetingRepo,
	}
}

func (uc *ListZoomRecordingsUseCase) Execute(ctx context.Context, zoomMeetingID string) ([]*dtos.ZoomRecordingDTO, error) {
	meeting, err := uc.meetingRepo.FindByID(ctx, zoomMeetingID)
	if err != nil {
		return nil, err
	}
	if meeting == nil {
		return nil, ErrZoomMeetingNotFound
	}

	recordings, err := uc.recordingRepo.FindByZoomMeetingID(ctx, zoomMeetingID)
	if err != nil {
		return nil, err
	}

	recordingDTOs := make([]*dtos.ZoomRecordingDTO, len(recordings))
	for i, recording := range recordings {
		var dto dtos.ZoomRecordingDTO
		dto.FromEntity(recording)
		recordingDTOs[i] = &dto
	}

	return recordingDTOs, nil
}

