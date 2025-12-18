package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
)

type ZoomMeetingRepository interface {
	Create(ctx context.Context, meeting *entities.ZoomMeeting) error
	FindByID(ctx context.Context, id string) (*entities.ZoomMeeting, error)
	FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) (*entities.ZoomMeeting, error)
	FindBySectionModuleID(ctx context.Context, sectionModuleID string) (*entities.ZoomMeeting, error)
	Update(ctx context.Context, meeting *entities.ZoomMeeting) error
	Delete(ctx context.Context, id string) error
}

