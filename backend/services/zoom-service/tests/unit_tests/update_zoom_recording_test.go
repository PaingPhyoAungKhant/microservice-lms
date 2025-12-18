package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateZoomRecording_Success(t *testing.T) {
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

	newRecordingType := "mov"
	newFileSize := int64(2048000)

	input := dtos.UpdateZoomRecordingInput{
		RecordingType: &newRecordingType,
		FileSize:      &newFileSize,
	}

	repo.On("FindByID", mock.Anything, "recording-id").Return(recording, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := usecases.NewUpdateZoomRecordingUseCase(repo)

	result, err := uc.Execute(context.Background(), "recording-id", input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, newRecordingType, *result.RecordingType)
	assert.Equal(t, newFileSize, *result.FileSize)
	repo.AssertExpectations(t)
}

func TestUpdateZoomRecording_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	newRecordingType := "mov"
	input := dtos.UpdateZoomRecordingInput{
		RecordingType: &newRecordingType,
	}

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, nil).Once()

	uc := usecases.NewUpdateZoomRecordingUseCase(repo)

	_, err := uc.Execute(context.Background(), "recording-id", input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomRecordingNotFound, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateZoomRecording_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomRecordingRepository)

	newRecordingType := "mov"
	input := dtos.UpdateZoomRecordingInput{
		RecordingType: &newRecordingType,
	}

	repo.On("FindByID", mock.Anything, "recording-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewUpdateZoomRecordingUseCase(repo)

	_, err := uc.Execute(context.Background(), "recording-id", input)

	require.Error(t, err)
	repo.AssertExpectations(t)
}

