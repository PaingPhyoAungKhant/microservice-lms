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

func TestGetZoomMeeting_Success(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	startTime := time.Now().UTC()
	duration := 60
	password := "pass123"

	meeting := entities.NewZoomMeeting(
		uuid.New().String(),
		"zoom-meeting-id-123",
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		&startTime,
		&duration,
		&password,
	)

	repo.On("FindByID", mock.Anything, "meeting-id").Return(meeting, nil).Once()

	uc := usecases.NewGetZoomMeetingUseCase(repo)

	result, err := uc.Execute(context.Background(), "meeting-id")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meeting.ID, result.ID)
	assert.Equal(t, meeting.Topic, result.Topic)
	repo.AssertExpectations(t)
}

func TestGetZoomMeeting_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, nil).Once()

	uc := usecases.NewGetZoomMeetingUseCase(repo)

	_, err := uc.Execute(context.Background(), "meeting-id")

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
	repo.AssertExpectations(t)
}

func TestGetZoomMeeting_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	repo.On("FindByID", mock.Anything, "meeting-id").Return(nil, errors.New("database error")).Once()

	uc := usecases.NewGetZoomMeetingUseCase(repo)

	_, err := uc.Execute(context.Background(), "meeting-id")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

