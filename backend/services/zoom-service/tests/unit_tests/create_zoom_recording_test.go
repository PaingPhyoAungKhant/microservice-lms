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

func TestCreateZoomRecording_Success(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()
	fileID := uuid.New().String()
	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour)
	fileSize := int64(1024000)
	recordingType := "mp4"

	meeting := entities.NewZoomMeeting(
		uuid.New().String(),
		"zoom-meeting-id-123",
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		nil,
		nil,
		nil,
	)

	input := dtos.CreateZoomRecordingInput{
		ZoomMeetingID:     meetingID,
		FileID:             fileID,
		RecordingType:      &recordingType,
		RecordingStartTime: &startTime,
		RecordingEndTime:   &endTime,
		FileSize:           &fileSize,
	}

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(meeting, nil).Once()
	recordingRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)

	result, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meetingID, result.ZoomMeetingID)
	assert.Equal(t, fileID, result.FileID)
	meetingRepo.AssertExpectations(t)
	recordingRepo.AssertExpectations(t)
}

func TestCreateZoomRecording_MeetingNotFound(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()
	fileID := uuid.New().String()
	recordingType := "mp4"

	input := dtos.CreateZoomRecordingInput{
		ZoomMeetingID: meetingID,
		FileID:         fileID,
		RecordingType:  &recordingType,
	}

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(nil, nil).Once()

	uc := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
	recordingRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateZoomRecording_RepositoryError(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()
	fileID := uuid.New().String()
	recordingType := "mp4"

	meeting := entities.NewZoomMeeting(
		uuid.New().String(),
		"zoom-meeting-id-123",
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		nil,
		nil,
		nil,
	)

	input := dtos.CreateZoomRecordingInput{
		ZoomMeetingID: meetingID,
		FileID:         fileID,
		RecordingType:  &recordingType,
	}

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(meeting, nil).Once()
	recordingRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()

	uc := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	meetingRepo.AssertExpectations(t)
	recordingRepo.AssertExpectations(t)
}

