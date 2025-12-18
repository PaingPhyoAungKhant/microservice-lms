package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type GetZoomMeetingUseCase struct {
	meetingRepo repositories.ZoomMeetingRepository
}

func NewGetZoomMeetingUseCase(meetingRepo repositories.ZoomMeetingRepository) *GetZoomMeetingUseCase {
	return &GetZoomMeetingUseCase{
		meetingRepo: meetingRepo,
	}
}

func (uc *GetZoomMeetingUseCase) Execute(ctx context.Context, id string) (*dtos.ZoomMeetingDTO, error) {
	meeting, err := uc.meetingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if meeting == nil {
		return nil, ErrZoomMeetingNotFound
	}

	var dto dtos.ZoomMeetingDTO
	dto.FromEntity(meeting)
	return &dto, nil
}

