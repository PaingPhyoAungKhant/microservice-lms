package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockZoomRecordingRepository struct {
	mock.Mock
}

func (m *MockZoomRecordingRepository) Create(ctx context.Context, recording *entities.ZoomRecording) error {
	args := m.Called(ctx, recording)
	return args.Error(0)
}

func (m *MockZoomRecordingRepository) FindByID(ctx context.Context, id string) (*entities.ZoomRecording, error) {
	args := m.Called(ctx, id)
	recording, _ := args.Get(0).(*entities.ZoomRecording)
	return recording, args.Error(1)
}

func (m *MockZoomRecordingRepository) FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) ([]*entities.ZoomRecording, error) {
	args := m.Called(ctx, zoomMeetingID)
	recordings, _ := args.Get(0).([]*entities.ZoomRecording)
	return recordings, args.Error(1)
}

func (m *MockZoomRecordingRepository) Update(ctx context.Context, recording *entities.ZoomRecording) error {
	args := m.Called(ctx, recording)
	return args.Error(0)
}

func (m *MockZoomRecordingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

