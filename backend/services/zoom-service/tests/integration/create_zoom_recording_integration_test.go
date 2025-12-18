package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateZoomRecording_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)
	recordingRepo := SetupZoomRecordingRepository(db)

	ctx := context.Background()

	meeting := entities.NewZoomMeeting(
		uuid.New().String(),
		"zoom-meeting-id-"+uuid.New().String(),
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		nil,
		nil,
		nil,
	)

	err := meetingRepo.Create(ctx, meeting)
	require.NoError(t, err)

	fileID := uuid.New().String()
	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour)
	fileSize := int64(1024000)
	recordingType := "mp4"

	input := dtos.CreateZoomRecordingInput{
		ZoomMeetingID:     meeting.ID,
		FileID:             fileID,
		RecordingType:      &recordingType,
		RecordingStartTime: &startTime,
		RecordingEndTime:   &endTime,
		FileSize:           &fileSize,
	}

	createUC := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)

	result, err := createUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meeting.ID, result.ZoomMeetingID)
	assert.Equal(t, fileID, result.FileID)
	assert.Equal(t, recordingType, *result.RecordingType)
}

func TestCreateZoomRecording_Integration_MeetingNotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)
	recordingRepo := SetupZoomRecordingRepository(db)

	fileID := uuid.New().String()
	recordingType := "mp4"

	input := dtos.CreateZoomRecordingInput{
		ZoomMeetingID: "00000000-0000-0000-0000-000000000000",
		FileID:         fileID,
		RecordingType:  &recordingType,
	}

	createUC := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)

	_, err := createUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

