package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var (
	ErrZoomMeetingAlreadyExists = errors.New("zoom meeting already exists for this section module")
)

type CreateZoomMeetingUseCase struct {
	meetingRepo repositories.ZoomMeetingRepository
	zoomClient  *zoom.ZoomClient
	publisher   messaging.Publisher
	logger      *logger.Logger
	userID      string
}

func NewCreateZoomMeetingUseCase(
	meetingRepo repositories.ZoomMeetingRepository,
	zoomClient *zoom.ZoomClient,
	publisher messaging.Publisher,
	logger *logger.Logger,
	userID string,
) *CreateZoomMeetingUseCase {
	return &CreateZoomMeetingUseCase{
		meetingRepo: meetingRepo,
		zoomClient:  zoomClient,
		publisher:   publisher,
		logger:      logger,
		userID:      userID,
	}
}

func (uc *CreateZoomMeetingUseCase) Execute(ctx context.Context, input dtos.CreateZoomMeetingInput) (*dtos.ZoomMeetingDTO, error) {
	existing, _ := uc.meetingRepo.FindBySectionModuleID(ctx, input.SectionModuleID)
	if existing != nil {
		return nil, ErrZoomMeetingAlreadyExists
	}

	var startTimeStr *string
	if input.StartTime != nil {
		formatted := input.StartTime.Format(time.RFC3339)
		startTimeStr = &formatted
	}

	zoomReq := zoom.CreateMeetingRequest{
		Topic:     input.Topic,
		Type:      2,
		StartTime: startTimeStr,
		Duration:  input.Duration,
		Password:  input.Password,
		Settings: &zoom.MeetingSettings{
			HostVideo:        true,
			ParticipantVideo: true,
			JoinBeforeHost:   false,
			MuteUponEntry:    false,
			WaitingRoom:      false,
		},
	}

	zoomResponse, err := uc.zoomClient.CreateMeeting(uc.userID, zoomReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create zoom meeting: %w", err)
	}

	var startTime *time.Time
	if zoomResponse.StartTime != "" {
		parsed, err := time.Parse(time.RFC3339, zoomResponse.StartTime)
		if err == nil {
			startTime = &parsed
		}
	}

	meeting := entities.NewZoomMeeting(
		input.SectionModuleID,
		string(zoomResponse.ID),
		zoomResponse.Topic,
		zoomResponse.JoinURL,
		zoomResponse.StartURL,
		startTime,
		&zoomResponse.Duration,
		&zoomResponse.Password,
	)

	if err := uc.meetingRepo.Create(ctx, meeting); err != nil {
		return nil, fmt.Errorf("failed to save zoom meeting: %w", err)
	}

	event := events.ZoomMeetingCreatedEvent{
		ZoomMeetingID:   meeting.ID,
		SectionModuleID: meeting.SectionModuleID,
		Topic:           meeting.Topic,
		JoinURL:         meeting.JoinURL,
		StartURL:        meeting.StartURL,
		CreatedAt:       meeting.CreatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeZoomMeetingCreated, event); err != nil {
			uc.logger.Error("Failed to publish zoom meeting created event", zap.Error(err))
		}
	}

	var dto dtos.ZoomMeetingDTO
	dto.FromEntity(meeting)
	return &dto, nil
}

