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

func TestCreateSectionModule_Success(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	section := &entities.CourseSection{
		ID:               sectionID,
		CourseOfferingID: "offering-123",
		Name:             "Introduction",
		Description:      "Introduction section",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	input := dtos.CreateSectionModuleInput{
		Name:        "Module 1",
		Description: "First module",
		ContentType: "zoom",
		Order:       1,
	}

	var createdModule *entities.SectionModule
	sectionRepo.On("FindByID", mock.Anything, sectionID).Return(section, nil).Once()
	moduleRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *entities.SectionModule) bool {
		createdModule = m
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeSectionModuleCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), sectionID, input)
	require.NoError(t, err)
	require.NotNil(t, createdModule)
	require.Equal(t, sectionID, createdModule.CourseSectionID)
	require.Equal(t, input.Name, createdModule.Name)
	require.Equal(t, input.Description, createdModule.Description)
	require.Equal(t, entities.ContentTypeZoom, createdModule.ContentType)
	require.Equal(t, entities.ContentStatusDraft, createdModule.ContentStatus)
	require.Equal(t, input.Order, createdModule.Order)
	require.Nil(t, createdModule.ContentID)

	assertDTOEqualSectionModule(t, dto, createdModule)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.SectionModuleCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdModule.ID, event.ID)
	assert.Equal(t, sectionID, event.CourseSectionID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, string(createdModule.ContentType), event.ContentType)

	moduleRepo.AssertExpectations(t)
	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateSectionModule_SectionNotFound(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	input := dtos.CreateSectionModuleInput{
		Name:        "Module 1",
		Description: "First module",
		ContentType: "zoom",
		Order:       1,
	}

	sectionRepo.On("FindByID", mock.Anything, sectionID).Return((*entities.CourseSection)(nil), nil).Once()

	uc := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), sectionID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseSectionNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	moduleRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	moduleRepo.AssertExpectations(t)
	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateSectionModule_RepositoryError(t *testing.T) {
	moduleRepo := new(mocks.MockSectionModuleRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	sectionID := "section-123"
	section := &entities.CourseSection{
		ID:               sectionID,
		CourseOfferingID: "offering-123",
		Name:             "Introduction",
		Description:      "Introduction section",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	input := dtos.CreateSectionModuleInput{
		Name:        "Module 1",
		Description: "First module",
		ContentType: "zoom",
		Order:       1,
	}

	sectionRepo.On("FindByID", mock.Anything, sectionID).Return(section, nil).Once()
	moduleRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), sectionID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	moduleRepo.AssertExpectations(t)
	sectionRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

