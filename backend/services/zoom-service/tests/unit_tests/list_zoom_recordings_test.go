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

func TestListZoomRecordings_Success(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()
	fileID1 := uuid.New().String()
	fileID2 := uuid.New().String()
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

	recording1 := entities.NewZoomRecording(
		meetingID,
		fileID1,
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	recording2 := entities.NewZoomRecording(
		meetingID,
		fileID2,
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(meeting, nil).Once()
	recordingRepo.On("FindByZoomMeetingID", mock.Anything, meetingID).Return([]*entities.ZoomRecording{recording1, recording2}, nil).Once()

	uc := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	result, err := uc.Execute(context.Background(), meetingID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
	meetingRepo.AssertExpectations(t)
	recordingRepo.AssertExpectations(t)
}

func TestListZoomRecordings_MeetingNotFound(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(nil, nil).Once()

	uc := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	_, err := uc.Execute(context.Background(), meetingID)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
	recordingRepo.AssertNotCalled(t, "FindByZoomMeetingID", mock.Anything, mock.Anything)
}

func TestListZoomRecordings_RepositoryError(t *testing.T) {
	meetingRepo := new(mocks.MockZoomMeetingRepository)
	recordingRepo := new(mocks.MockZoomRecordingRepository)

	meetingID := uuid.New().String()

	meetingRepo.On("FindByID", mock.Anything, meetingID).Return(nil, errors.New("database error")).Once()

	uc := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	_, err := uc.Execute(context.Background(), meetingID)

	require.Error(t, err)
	meetingRepo.AssertExpectations(t)
}

