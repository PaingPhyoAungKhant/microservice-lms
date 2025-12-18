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

func TestCreateSectionModule_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	section := entities.NewCourseSection(offering.ID, "Introduction", "Introduction section", 1)
	err = sectionRepo.Create(ctx, section)
	require.NoError(t, err)

	input := dtos.CreateSectionModuleInput{
		Name:        "Module 1",
		Description: "First module",
		ContentType: "zoom",
		Order:       1,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	createModuleUC := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	result, err := createModuleUC.Execute(ctx, section.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.ContentType, result.ContentType)
	assert.Equal(t, "draft", result.ContentStatus)
	assert.Equal(t, input.Order, result.Order)
	assert.Nil(t, result.ContentID)

	createdModule, err := moduleRepo.FindByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, createdModule)
	assert.Equal(t, input.Name, createdModule.Name)
	assert.Equal(t, section.ID, createdModule.CourseSectionID)

	publisher.AssertExpectations(t)
}

func TestCreateSectionModule_Integration_SectionNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.CreateSectionModuleInput{
		Name:        "Module 1",
		Description: "First module",
		ContentType: "zoom",
		Order:       1,
	}

	createModuleUC := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	nonExistentSectionID := uuid.New().String()
	_, err = createModuleUC.Execute(ctx, nonExistentSectionID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseSectionNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateSectionModule_Integration_MultipleModulesWithOrdering(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	section := entities.NewCourseSection(offering.ID, "Introduction", "Introduction section", 1)
	err = sectionRepo.Create(ctx, section)
	require.NoError(t, err)

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)

	createModuleUC := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, publisher, logger)

	modules := []dtos.CreateSectionModuleInput{
		{Name: "Module 1", Description: "First module", ContentType: "zoom", Order: 1},
		{Name: "Module 2", Description: "Second module", ContentType: "zoom", Order: 2},
		{Name: "Module 3", Description: "Third module", ContentType: "zoom", Order: 3},
	}

	for _, input := range modules {
		result, err := createModuleUC.Execute(ctx, section.ID, input)
		require.NoError(t, err)
		assert.Equal(t, input.Order, result.Order)
	}

	createdModules, err := moduleRepo.FindBySectionID(ctx, section.ID)
	require.NoError(t, err)
	require.Len(t, createdModules, 3)
	assert.Equal(t, 1, createdModules[0].Order)
	assert.Equal(t, 2, createdModules[1].Order)
	assert.Equal(t, 3, createdModules[2].Order)

	publisher.AssertExpectations(t)
}

