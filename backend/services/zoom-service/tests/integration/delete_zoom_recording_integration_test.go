package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteZoomRecording_Integration_Success(t *testing.T) {
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

	deleteUC := usecases.NewDeleteZoomRecordingUseCase(recordingRepo)

	err = deleteUC.Execute(ctx, recording.ID)
	require.NoError(t, err)

	deletedRecording, err := recordingRepo.FindByID(ctx, recording.ID)
	require.Error(t, err)
	assert.Nil(t, deletedRecording)
}

func TestDeleteZoomRecording_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	recordingRepo := SetupZoomRecordingRepository(db)

	deleteUC := usecases.NewDeleteZoomRecordingUseCase(recordingRepo)

	err := deleteUC.Execute(context.Background(), "00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

