package unit_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/tests/mocks"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateZoomMeeting_AlreadyExists(t *testing.T) {
	repo := new(mocks.MockZoomMeetingRepository)
	zoomClient := &zoom.ZoomClient{}
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	moduleID := uuid.New().String()
	existingMeeting := entities.NewZoomMeeting(
		moduleID,
		"zoom-meeting-id-123",
		"Existing Meeting",
		"https://zoom.us/j/123",
		"https://zoom.us/s/123",
		nil,
		nil,
		nil,
	)

	input := dtos.CreateZoomMeetingInput{
		SectionModuleID: moduleID,
		Topic:           "New Meeting",
	}

	repo.On("FindBySectionModuleID", mock.Anything, moduleID).Return(existingMeeting, nil).Once()

	uc := usecases.NewCreateZoomMeetingUseCase(repo, zoomClient, publisher, logger, "user-123")

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrZoomMeetingAlreadyExists, err)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

