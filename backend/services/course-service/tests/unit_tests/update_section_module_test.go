package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateSectionModule_Success(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	moduleID := "module-123"
	contentID := "content-123"
	existingModule := &entities.SectionModule{
		ID:              moduleID,
		CourseSectionID: "section-123",
		ContentID:       &contentID,
		Name:            "Module 1",
		Description:     "Old description",
		ContentType:     entities.ContentTypeZoom,
		ContentStatus:   entities.ContentStatusDraft,
		Order:           1,
	}

	input := dtos.UpdateSectionModuleInput{
		Name:        "Updated Module",
		Description: "Updated description",
		Order:       2,
	}

	moduleRepo.On("FindByID", mock.Anything, moduleID).Return(existingModule, nil).Once()
	moduleRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *entities.SectionModule) bool {
		return m.ID == moduleID && m.Name == input.Name && m.Description == input.Description
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeSectionModuleUpdated, mock.Anything).Return(nil).Once()

	uc := usecases.NewUpdateSectionModuleUseCase(moduleRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), moduleID, input)
	require.NoError(t, err)
	require.Equal(t, input.Name, existingModule.Name)
	require.Equal(t, input.Description, existingModule.Description)
	require.Equal(t, input.Order, existingModule.Order)

	assertDTOEqualSectionModule(t, dto, existingModule)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.SectionModuleUpdatedEvent)
	require.True(t, ok)
	assert.Equal(t, existingModule.ID, event.ID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, input.Order, event.Order)

	moduleRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateSectionModule_ModuleNotFound(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	moduleID := "module-123"
	input := dtos.UpdateSectionModuleInput{
		Name:        "Updated Module",
		Description: "Updated description",
		Order:       2,
	}

	moduleRepo.On("FindByID", mock.Anything, moduleID).Return((*entities.SectionModule)(nil), nil).Once()

	uc := usecases.NewUpdateSectionModuleUseCase(moduleRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), moduleID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrSectionModuleNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	moduleRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	moduleRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateSectionModule_RepositoryError(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	moduleID := "module-123"
	existingModule := &entities.SectionModule{
		ID:              moduleID,
		CourseSectionID: "section-123",
		Name:            "Module 1",
		Description:     "Old description",
		ContentType:     entities.ContentTypeZoom,
		ContentStatus:   entities.ContentStatusDraft,
		Order:           1,
	}

	input := dtos.UpdateSectionModuleInput{
		Name:        "Updated Module",
		Description: "Updated description",
		Order:       2,
	}

	moduleRepo.On("FindByID", mock.Anything, moduleID).Return(existingModule, nil).Once()
	moduleRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewUpdateSectionModuleUseCase(moduleRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), moduleID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	moduleRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

