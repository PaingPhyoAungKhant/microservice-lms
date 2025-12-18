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

func TestDeleteZoomRecording_Success(t *testing.T) {
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
	repo.On("Delete", mock.Anything, "recording-id").Return(nil).Once()

	uc := usecases.NewDeleteZoomRecordingUseCase(repo)

	err := uc.Execute(context.Background(), "recording-id")

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteZoomRecording_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, nil).Once()

	uc := usecases.NewDeleteZoomRecordingUseCase(repo)

	err := uc.Execute(context.Background(), "recording-id")

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomRecordingNotFound, err)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestDeleteZoomRecording_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewDeleteZoomRecordingUseCase(repo)

	err := uc.Execute(context.Background(), "recording-id")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteZoomRecording_DeleteError(t *testing.T) {
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
	repo.On("Delete", mock.Anything, "recording-id").Return(errors.New("delete failed")).Once()

	uc := usecases.NewDeleteZoomRecordingUseCase(repo)

	err := uc.Execute(context.Background(), "recording-id")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

