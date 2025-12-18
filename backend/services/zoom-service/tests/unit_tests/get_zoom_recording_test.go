package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetZoomRecording_Success(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour)
	fileSize := int64(1024000)
	recordingType := "mp4"

	recording := entities.NewZoomRecording(
		uuid.New().String(),
		uuid.New().String(),
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	repo.On("FindByID", mock.Anything, "recording-id").Return(recording, nil).Once()

	uc := usecases.NewGetZoomRecordingUseCase(repo)

	result, err := uc.Execute(context.Background(), "recording-id")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, recording.ID, result.ID)
	assert.Equal(t, recording.ZoomMeetingID, result.ZoomMeetingID)
	repo.AssertExpectations(t)
}

func TestGetZoomRecording_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, nil).Once()

	uc := usecases.NewGetZoomRecordingUseCase(repo)

	_, err := uc.Execute(context.Background(), "recording-id")

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomRecordingNotFound, err)
	repo.AssertExpectations(t)
}

func TestGetZoomRecording_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewGetZoomRecordingUseCase(repo)

	_, err := uc.Execute(context.Background(), "recording-id")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

