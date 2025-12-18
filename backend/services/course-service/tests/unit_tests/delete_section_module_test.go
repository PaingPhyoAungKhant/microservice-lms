package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteSectionModule_FindByIDError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "module-1").Return((*entities.SectionModule)(nil), expectedErr).Once()

	uc := usecases.NewDeleteSectionModuleUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteSectionModuleInput{ModuleID: "module-1"})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteSectionModule_ModuleNotFound(t *testing.T) {
	repo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "module-1").Return((*entities.SectionModule)(nil), nil).Once()

	uc := usecases.NewDeleteSectionModuleUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteSectionModuleInput{ModuleID: "module-1"})
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrSectionModuleNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteSectionModule_Success(t *testing.T) {
	module := &entities.SectionModule{
		ID:              "module-1",
		CourseSectionID: "section-1",
		Name:            "Test Module",
		Description:     "Test Description",
		ContentType:     entities.ContentTypeZoom,
		ContentStatus:   entities.ContentStatusDraft,
		Order:           1,
	}

	repo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)

	repo.On("FindByID", mock.Anything, module.ID).Return(module, nil).Once()
	repo.On("Delete", mock.Anything, module.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeSectionModuleDeleted, mock.Anything).Return(nil).Once()

	uc := usecases.NewDeleteSectionModuleUseCase(repo, publisher, logger.NewNop())

	out, err := uc.Execute(context.Background(), usecases.DeleteSectionModuleInput{ModuleID: module.ID})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotEmpty(t, out.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.SectionModuleDeletedEvent)
	require.True(t, ok)
	assert.Equal(t, module.ID, event.ID)
	assert.WithinDuration(t, time.Now(), event.DeletedAt, 2*time.Second)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

