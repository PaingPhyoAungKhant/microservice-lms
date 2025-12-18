package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	sharedMocks "github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourseSection_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	input := dtos.CreateCourseSectionInput{
		Name:        "Introduction",
		Description: "Introduction to the course",
		Order:       1,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	createSectionUC := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	result, err := createSectionUC.Execute(ctx, offering.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.Order, result.Order)
	assert.Equal(t, "draft", result.Status)

	createdSection, err := sectionRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, createdSection)
	assert.Equal(t, input.Name, createdSection.Name)
	assert.Equal(t, offering.ID, createdSection.CourseOfferingID)

	publisher.AssertExpectations(t)
}

func TestCreateCourseSection_Integration_OfferingNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	offeringRepo := SetupCourseOfferingRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.CreateCourseSectionInput{
		Name:        "Introduction",
		Description: "Introduction to the course",
		Order:       1,
	}

	createSectionUC := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	nonExistentOfferingID := uuid.New().String()
	_, err = createSectionUC.Execute(ctx, nonExistentOfferingID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseOfferingNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateCourseSection_Integration_MultipleSectionsWithOrdering(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)

	createSectionUC := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, publisher, logger)

	sections := []dtos.CreateCourseSectionInput{
		{Name: "Section 1", Description: "First section", Order: 1},
		{Name: "Section 2", Description: "Second section", Order: 2},
		{Name: "Section 3", Description: "Third section", Order: 3},
	}

	for _, input := range sections {
		result, err := createSectionUC.Execute(ctx, offering.ID, input)
		require.NoError(t, err)
		assert.Equal(t, input.Order, result.Order)
	}

	createdSections, err := sectionRepo.FindByOfferingID(ctx, offering.ID)
	require.NoError(t, err)
	require.Len(t, createdSections, 3)
	assert.Equal(t, 1, createdSections[0].Order)
	assert.Equal(t, 2, createdSections[1].Order)
	assert.Equal(t, 3, createdSections[2].Order)

	publisher.AssertExpectations(t)
}

