package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteZoomMeeting_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)
	zoomClient := &zoom.ZoomClient{}

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, nil).Once()

	uc := usecases.NewDeleteZoomMeetingUseCase(repo, zoomClient)

	err := uc.Execute(context.Background(), "meeting-id")

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
}

func TestDeleteZoomMeeting_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)
	zoomClient := &zoom.ZoomClient{}

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewDeleteZoomMeetingUseCase(repo, zoomClient)

	err := uc.Execute(context.Background(), "meeting-id")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

