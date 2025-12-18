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

func TestCreateCourseSection_Success(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     "course-123",
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	input := dtos.CreateCourseSectionInput{
		Name:        "Introduction",
		Description: "Introduction to the course",
		Order:       1,
	}

	var createdSection *entities.CourseSection
	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(offering, nil).Once()
	sectionRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *entities.CourseSection) bool {
		createdSection = s
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeCourseSectionCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	dto, err := uc.Execute(context.Background(), offeringID, input)
	require.NoError(t, err)
	require.NotNil(t, createdSection)
	require.Equal(t, offeringID, createdSection.CourseOfferingID)
	require.Equal(t, input.Name, createdSection.Name)
	require.Equal(t, input.Description, createdSection.Description)
	require.Equal(t, input.Order, createdSection.Order)
	require.Equal(t, entities.SectionStatusDraft, createdSection.Status)

	assertDTOEqualCourseSection(t, dto, createdSection)

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.CourseSectionCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdSection.ID, event.ID)
	assert.Equal(t, offeringID, event.CourseOfferingID)
	assert.Equal(t, input.Name, event.Name)
	assert.Equal(t, input.Order, event.Order)

	sectionRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourseSection_OfferingNotFound(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	input := dtos.CreateCourseSectionInput{
		Name:        "Introduction",
		Description: "Introduction to the course",
		Order:       1,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return((*entities.CourseOffering)(nil), nil).Once()

	uc := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseOfferingNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	sectionRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	sectionRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateCourseSection_RepositoryError(t *testing.T) {
	sectionRepo := new(mocks.MockCourseSectionRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     "course-123",
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	input := dtos.CreateCourseSectionInput{
		Name:        "Introduction",
		Description: "Introduction to the course",
		Order:       1,
	}

	offeringRepo.On("FindByID", mock.Anything, offeringID).Return(offering, nil).Once()
	sectionRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	uc := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	_, err := uc.Execute(context.Background(), offeringID, input)
	require.Error(t, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	sectionRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

