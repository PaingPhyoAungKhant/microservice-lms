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

func TestListZoomRecordings_Integration_Success(t *testing.T) {
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

	fileID1 := uuid.New().String()
	fileID2 := uuid.New().String()
	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Hour)
	fileSize := int64(1024000)
	recordingType := "mp4"

	recording1 := entities.NewZoomRecording(
		meeting.ID,
		fileID1,
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	recording2 := entities.NewZoomRecording(
		meeting.ID,
		fileID2,
		&recordingType,
		&startTime,
		&endTime,
		&fileSize,
	)

	err = recordingRepo.Create(ctx, recording1)
	require.NoError(t, err)
	err = recordingRepo.Create(ctx, recording2)
	require.NoError(t, err)

	listUC := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	result, err := listUC.Execute(ctx, meeting.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestListZoomRecordings_Integration_MeetingNotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)
	recordingRepo := SetupZoomRecordingRepository(db)

	listUC := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	_, err := listUC.Execute(context.Background(), "00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

