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

func TestUpdateZoomRecording_Integration_Success(t *testing.T) {
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

	recording := entities.NewZoomRecording(
		meeting.ID,
		fileID,
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	err = recordingRepo.Create(ctx, recording)
	require.NoError(t, err)

	newRecordingType := "mov"
	newFileSize := int64(2048000)

	input := dtos.UpdateZoomRecordingInput{
		RecordingType: &newRecordingType,
		FileSize:      &newFileSize,
	}

	updateUC := usecases.NewUpdateZoomRecordingUseCase(recordingRepo)

	result, err := updateUC.Execute(ctx, recording.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, newRecordingType, *result.RecordingType)
	assert.Equal(t, newFileSize, *result.FileSize)

	updatedRecording, err := recordingRepo.FindByID(ctx, recording.ID)
	require.NoError(t, err)
	assert.Equal(t, newRecordingType, *updatedRecording.RecordingType)
	assert.Equal(t, newFileSize, *updatedRecording.FileSize)
}

func TestUpdateZoomRecording_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	recordingRepo := SetupZoomRecordingRepository(db)

	newRecordingType := "mov"
	input := dtos.UpdateZoomRecordingInput{
		RecordingType: &newRecordingType,
	}

	updateUC := usecases.NewUpdateZoomRecordingUseCase(recordingRepo)

	_, err := updateUC.Execute(context.Background(), "00000000-0000-0000-0000-000000000000", input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

