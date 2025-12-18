package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
)

var (
	ErrZoomMeetingNotFound = errors.New("zoom meeting not found")
)

type UpdateZoomMeetingUseCase struct {
	meetingRepo repositories.ZoomMeetingRepository
	zoomClient  *zoom.ZoomClient
}

func NewUpdateZoomMeetingUseCase(
	meetingRepo repositories.ZoomMeetingRepository,
	zoomClient *zoom.ZoomClient,
) *UpdateZoomMeetingUseCase {
	return &UpdateZoomMeetingUseCase{
		meetingRepo: meetingRepo,
		zoomClient:  zoomClient,
	}
}

func (uc *UpdateZoomMeetingUseCase) Execute(ctx context.Context, id string, input dtos.UpdateZoomMeetingInput) (*dtos.ZoomMeetingDTO, error) {
	meeting, err := uc.meetingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if meeting == nil {
		return nil, ErrZoomMeetingNotFound
	}

	var startTimeStr *string
	if input.StartTime != nil {
		formatted := input.StartTime.Format(time.RFC3339)
		startTimeStr = &formatted
	}

	zoomReq := zoom.UpdateMeetingRequest{
		Topic:     input.Topic,
		StartTime: startTimeStr,
		Duration:  input.Duration,
		Password:  input.Password,
	}

	if err := uc.zoomClient.UpdateMeeting(meeting.ZoomMeetingID, zoomReq); err != nil {
		return nil, fmt.Errorf("failed to update zoom meeting: %w", err)
	}

	meeting.Update(input.Topic, input.StartTime, input.Duration, input.Password)

	zoomResponse, err := uc.zoomClient.GetMeeting(meeting.ZoomMeetingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated zoom meeting: %w", err)
	}

	var startTime *time.Time
	if zoomResponse.StartTime != "" {
		parsed, err := time.Parse(time.RFC3339, zoomResponse.StartTime)
		if err == nil {
			startTime = &parsed
		}
	}

	meeting.UpdateURLs(zoomResponse.JoinURL, zoomResponse.StartURL)
	if startTime != nil {
		meeting.StartTime = startTime
	}

	if err := uc.meetingRepo.Update(ctx, meeting); err != nil {
		return nil, fmt.Errorf("failed to update zoom meeting in database: %w", err)
	}

	var dto dtos.ZoomMeetingDTO
	dto.FromEntity(meeting)
	return &dto, nil
}

