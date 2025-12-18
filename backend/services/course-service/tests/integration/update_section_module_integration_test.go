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

func TestUpdateSectionModule_Integration_Success(t *testing.T) {
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

	module := entities.NewSectionModule(section.ID, "Module 1", "First module", entities.ContentTypeZoom, 1)
	err = moduleRepo.Create(ctx, module)
	require.NoError(t, err)

	input := dtos.UpdateSectionModuleInput{
		Name:        "Updated Module",
		Description: "Updated description",
		Order:       2,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	updateModuleUC := usecases.NewUpdateSectionModuleUseCase(moduleRepo, publisher, logger)

	result, err := updateModuleUC.Execute(ctx, module.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.Order, result.Order)

	updatedModule, err := moduleRepo.FindByID(ctx, module.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedModule)
	assert.Equal(t, input.Name, updatedModule.Name)
	assert.Equal(t, input.Description, updatedModule.Description)

	publisher.AssertExpectations(t)
}

func TestUpdateSectionModule_Integration_ModuleNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	moduleRepo := SetupSectionModuleRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.UpdateSectionModuleInput{
		Name:        "Updated Module",
		Description: "Updated description",
		Order:       2,
	}

	updateModuleUC := usecases.NewUpdateSectionModuleUseCase(moduleRepo, publisher, logger)

	nonExistentModuleID := uuid.New().String()
	_, err = updateModuleUC.Execute(ctx, nonExistentModuleID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrSectionModuleNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

