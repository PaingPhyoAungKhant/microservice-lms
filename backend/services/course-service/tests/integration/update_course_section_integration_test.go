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

func TestUpdateCourseSection_Integration_Success(t *testing.T) {
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

	section := entities.NewCourseSection(offering.ID, "Introduction", "Introduction section", 1)
	err = sectionRepo.Create(ctx, section)
	require.NoError(t, err)

	input := dtos.UpdateCourseSectionInput{
		Name:        "Advanced Topics",
		Description: "Updated description",
		Order:       2,
	}

	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	updateSectionUC := usecases.NewUpdateCourseSectionUseCase(sectionRepo, publisher, logger)

	result, err := updateSectionUC.Execute(ctx, section.ID, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.Order, result.Order)

	updatedSection, err := sectionRepo.FindByID(ctx, section.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedSection)
	assert.Equal(t, input.Name, updatedSection.Name)
	assert.Equal(t, input.Description, updatedSection.Description)

	publisher.AssertExpectations(t)
}

func TestUpdateCourseSection_Integration_SectionNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"course_section", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	sectionRepo := SetupCourseSectionRepository(db)
	publisher := new(sharedMocks.MockPublisher)
	logger := logger.NewNop()

	ctx := context.Background()

	input := dtos.UpdateCourseSectionInput{
		Name:        "Advanced Topics",
		Description: "Updated description",
		Order:       2,
	}

	updateSectionUC := usecases.NewUpdateCourseSectionUseCase(sectionRepo, publisher, logger)

	nonExistentSectionID := uuid.New().String()
	_, err = updateSectionUC.Execute(ctx, nonExistentSectionID, input)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseSectionNotFound, err)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

