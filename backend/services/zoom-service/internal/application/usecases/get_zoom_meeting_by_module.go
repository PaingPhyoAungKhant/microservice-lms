package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type GetZoomMeetingByModuleUseCase struct {
	meetingRepo repositories.ZoomMeetingRepository
}

func NewGetZoomMeetingByModuleUseCase(meetingRepo repositories.ZoomMeetingRepository) *GetZoomMeetingByModuleUseCase {
	return &GetZoomMeetingByModuleUseCase{
		meetingRepo: meetingRepo,
	}
}

func (uc *GetZoomMeetingByModuleUseCase) Execute(ctx context.Context, sectionModuleID string) (*dtos.ZoomMeetingDTO, error) {
	meeting, err := uc.meetingRepo.FindBySectionModuleID(ctx, sectionModuleID)
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

