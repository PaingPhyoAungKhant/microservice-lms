package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateZoomMeeting_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)
	zoomClient := &zoom.ZoomClient{}

	input := dtos.UpdateZoomMeetingInput{
		Topic: "Updated Topic",
	}

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, nil).Once()

	uc := usecases.NewUpdateZoomMeetingUseCase(repo, zoomClient)

	_, err := uc.Execute(context.Background(), "meeting-id", input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
}

func TestUpdateZoomMeeting_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)
	zoomClient := &zoom.ZoomClient{}

	input := dtos.UpdateZoomMeetingInput{
		Topic: "Updated Topic",
	}

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewUpdateZoomMeetingUseCase(repo, zoomClient)

	_, err := uc.Execute(context.Background(), "meeting-id", input)

	require.Error(t, err)
	repo.AssertExpectations(t)
}

