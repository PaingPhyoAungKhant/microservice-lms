package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCourseWithDetails_Integration_Success(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	instructorID := uuid.New().String()
	instructor := entities.NewCourseOfferingInstructor(offering.ID, instructorID, "instructor_user")
	err = instructorRepo.Create(ctx, instructor)
	require.NoError(t, err)

	section := entities.NewCourseSection(offering.ID, "Introduction", "Introduction section", 1)
	err = sectionRepo.Create(ctx, section)
	require.NoError(t, err)

	module := entities.NewSectionModule(section.ID, "Module 1", "First module", entities.ContentTypeZoom, 1)
	err = moduleRepo.Create(ctx, module)
	require.NoError(t, err)

	getCourseWithDetailsUC := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := getCourseWithDetailsUC.Execute(ctx, course.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, course.ID, result.Course.ID)
	require.Len(t, result.Offerings, 1)
	assert.Equal(t, offering.ID, result.Offerings[0].ID)
	require.Len(t, result.Offerings[0].Instructors, 1)
	assert.Equal(t, instructorID, result.Offerings[0].Instructors[0].InstructorID)
	require.Len(t, result.Offerings[0].Sections, 1)
	assert.Equal(t, section.ID, result.Offerings[0].Sections[0].ID)
	require.Len(t, result.Offerings[0].Sections[0].Modules, 1)
	assert.Equal(t, module.ID, result.Offerings[0].Sections[0].Modules[0].ID)
}

func TestGetCourseWithDetails_Integration_CourseNotFound(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)

	getCourseWithDetailsUC := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	nonExistentCourseID := uuid.New().String()
	_, err = getCourseWithDetailsUC.Execute(context.Background(), nonExistentCourseID)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrCourseNotFound, err)
}

func TestGetCourseWithDetails_Integration_CourseWithNoOfferings(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	getCourseWithDetailsUC := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := getCourseWithDetailsUC.Execute(ctx, course.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, course.ID, result.Course.ID)
	require.Len(t, result.Offerings, 0)
}

func TestGetCourseWithDetails_Integration_OfferingWithNoSections(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)

	ctx := context.Background()

	course := entities.NewCourse("Test Course", "Test Description", nil)
	err = courseRepo.Create(ctx, course)
	require.NoError(t, err)

	offering := entities.NewCourseOffering(course.ID, "Spring 2024", "Spring offering", entities.OfferingTypeOnline, nil, nil, 0.0)
	err = offeringRepo.Create(ctx, offering)
	require.NoError(t, err)

	getCourseWithDetailsUC := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := getCourseWithDetailsUC.Execute(ctx, course.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Offerings, 1)
	require.Len(t, result.Offerings[0].Instructors, 0)
	require.Len(t, result.Offerings[0].Sections, 0)
}

func TestGetCourseWithDetails_Integration_SectionWithNoModules(t *testing.T) {
	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, sharedIntegration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"section_module", "course_section", "course_offering_instructor", "course_offering", "course"},
	})
	require.NoError(t, err)
	defer cleanup()

	courseRepo := SetupCourseRepository(db)
	offeringRepo := SetupCourseOfferingRepository(db)
	instructorRepo := SetupCourseOfferingInstructorRepository(db)
	sectionRepo := SetupCourseSectionRepository(db)
	moduleRepo := SetupSectionModuleRepository(db)

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

	getCourseWithDetailsUC := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := getCourseWithDetailsUC.Execute(ctx, course.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Offerings, 1)
	require.Len(t, result.Offerings[0].Sections, 1)
	require.Len(t, result.Offerings[0].Sections[0].Modules, 0)
}

