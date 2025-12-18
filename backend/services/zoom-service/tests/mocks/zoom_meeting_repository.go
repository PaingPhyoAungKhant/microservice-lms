package mocks

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockZoomMeetingRepository struct {
	mock.Mock
}

func (m *MockZoomMeetingRepository) Create(ctx context.Context, meeting *entities.ZoomMeeting) error {
	args := m.Called(ctx, meeting)
	return args.Error(0)
}

func (m *MockZoomMeetingRepository) FindByID(ctx context.Context, id string) (*entities.ZoomMeeting, error) {
	args := m.Called(ctx, id)
	meeting, _ := args.Get(0).(*entities.ZoomMeeting)
	return meeting, args.Error(1)
}

func (m *MockZoomMeetingRepository) FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) (*entities.ZoomMeeting, error) {
	args := m.Called(ctx, zoomMeetingID)
	meeting, _ := args.Get(0).(*entities.ZoomMeeting)
	return meeting, args.Error(1)
}

func (m *MockZoomMeetingRepository) FindBySectionModuleID(ctx context.Context, sectionModuleID string) (*entities.ZoomMeeting, error) {
	args := m.Called(ctx, sectionModuleID)
	meeting, _ := args.Get(0).(*entities.ZoomMeeting)
	return meeting, args.Error(1)
}

func (m *MockZoomMeetingRepository) Update(ctx context.Context, meeting *entities.ZoomMeeting) error {
	args := m.Called(ctx, meeting)
	return args.Error(0)
}

func (m *MockZoomMeetingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

