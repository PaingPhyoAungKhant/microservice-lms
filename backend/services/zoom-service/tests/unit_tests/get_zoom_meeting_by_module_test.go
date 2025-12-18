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

func TestGetZoomMeetingByModule_Success(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	startTime := time.Now().UTC()
	duration := 60
	password := "pass123"
	moduleID := uuid.New().String()

	meeting := entities.NewZoomMeeting(
		moduleID,
		"zoom-meeting-id-123",
		"Test Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		&startTime,
		&duration,
		&password,
	)

	repo.On("FindBySectionModuleID", mock.Anything, moduleID).Return(meeting, nil).Once()

	uc := usecases.NewGetZoomMeetingByModuleUseCase(repo)

	result, err := uc.Execute(context.Background(), moduleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, meeting.ID, result.ID)
	assert.Equal(t, meeting.SectionModuleID, result.SectionModuleID)
	repo.AssertExpectations(t)
}

func TestGetZoomMeetingByModule_NotFound(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	moduleID := uuid.New().String()

	repo.On("FindBySectionModuleID", mock.Anything, moduleID).Return(nil, nil).Once()

	uc := usecases.NewGetZoomMeetingByModuleUseCase(repo)

	_, err := uc.Execute(context.Background(), moduleID)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingNotFound, err)
	repo.AssertExpectations(t)
}

func TestGetZoomMeetingByModule_RepositoryError(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)

	moduleID := uuid.New().String()

	repo.On("FindBySectionModuleID", mock.Anything, moduleID).Return(nil, errors.New("database error")).Once()

	uc := usecases.NewGetZoomMeetingByModuleUseCase(repo)

	_, err := uc.Execute(context.Background(), moduleID)

	require.Error(t, err)
	repo.AssertExpectations(t)
}

