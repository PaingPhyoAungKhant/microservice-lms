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

func TestGetZoomMeetingByModule_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)

	ctx := context.Background()

	moduleID := uuid.New().String()
	startTime := time.Now().UTC()
	duration := 60
	password := "pass123"

	meeting := entities.NewZoomMeeting(
		moduleID,
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

	getUC := usecases.NewGetZoomMeetingByModuleUseCase(meetingRepo)

	result, err := getUC.Execute(ctx, moduleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meeting.ID, result.ID)
	assert.Equal(t, moduleID, result.SectionModuleID)
}

func TestGetZoomMeetingByModule_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	meetingRepo := SetupZoomMeetingRepository(db)

	getUC := usecases.NewGetZoomMeetingByModuleUseCase(meetingRepo)

	_, err := getUC.Execute(context.Background(), uuid.New().String())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

