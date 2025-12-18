package usecases

import (
	"context"
	"fmt"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
)

type DeleteZoomMeetingUseCase struct {
	meetingRepo repositories.ZoomMeetingRepository
	zoomClient  *zoom.ZoomClient
}

func NewDeleteZoomMeetingUseCase(
	meetingRepo repositories.ZoomMeetingRepository,
	zoomClient *zoom.ZoomClient,
) *DeleteZoomMeetingUseCase {
	return &DeleteZoomMeetingUseCase{
		meetingRepo: meetingRepo,
		zoomClient:  zoomClient,
	}
}

func (uc *DeleteZoomMeetingUseCase) Execute(ctx context.Context, id string) error {
	meeting, err := uc.meetingRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if meeting == nil {
		return ErrZoomMeetingNotFound
	}

	if err := uc.zoomClient.DeleteMeeting(meeting.ZoomMeetingID); err != nil {
		return fmt.Errorf("failed to delete zoom meeting: %w", err)
	}

	if err := uc.meetingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete zoom meeting from database: %w", err)
	}

	return nil
}

