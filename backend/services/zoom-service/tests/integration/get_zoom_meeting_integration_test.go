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

func TestGetZoomMeeting_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)

	ctx := context.Background()

	startTime := time.Now().UTC()
	duration := 60
	password := "pass123"

	meeting := entities.NewZoomMeeting(
		uuid.New().String(),
		"zoom-meeting-id-"+uuid.New().String(),
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		&startTime,
		&duration,
		&password,
	)

	err := meetingRepo.Create(ctx, meeting)
	require.NoError(t, err)

	getUC := usecases.NewGetZoomMeetingUseCase(meetingRepo)

	result, err := getUC.Execute(ctx, meeting.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meeting.ID, result.ID)
	assert.Equal(t, meeting.Topic, result.Topic)
	assert.Equal(t, meeting.JoinURL, result.JoinURL)
}

func TestGetZoomMeeting_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)

	getUC := usecases.NewGetZoomMeetingUseCase(meetingRepo)

	_, err := getUC.Execute(context.Background(), "00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

