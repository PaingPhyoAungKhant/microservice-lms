package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
)

type ZoomRecordingRepository interface {
	Create(ctx context.Context, recording *entities.ZoomRecording) error
	FindByID(ctx context.Context, id string) (*entities.ZoomRecording, error)
	FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) ([]*entities.ZoomRecording, error)
	Update(ctx context.Context, recording *entities.ZoomRecording) error
	Delete(ctx context.Context, id string) error
}

