package unit_test

import (
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func assertDTOEqualCategory(t *testing.T, dto *dtos.CategoryDTO, category *entities.Category) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, category.ID, dto.ID)
	require.Equal(t, category.Name, dto.Name)
	require.Equal(t, category.Description, dto.Description)
}

func assertDTOEqualCourse(t *testing.T, dto *dtos.CourseDTO, course *entities.Course) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, course.ID, dto.ID)
	require.Equal(t, course.Name, dto.Name)
	require.Equal(t, course.Description, dto.Description)
	if course.ThumbnailID != nil {
		require.Equal(t, course.ThumbnailID, dto.ThumbnailID)
	}
}

func assertDTOEqualCourseOffering(t *testing.T, dto *dtos.CourseOfferingDTO, offering *entities.CourseOffering) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, offering.ID, dto.ID)
	require.Equal(t, offering.CourseID, dto.CourseID)
	require.Equal(t, offering.Name, dto.Name)
	require.Equal(t, offering.Description, dto.Description)
	require.Equal(t, string(offering.OfferingType), dto.OfferingType)
	require.Equal(t, string(offering.Status), dto.Status)
	if offering.Duration != nil {
		require.Equal(t, offering.Duration, dto.Duration)
	}
	if offering.ClassTime != nil {
		require.Equal(t, offering.ClassTime, dto.ClassTime)
	}
	require.Equal(t, offering.EnrollmentCost, dto.EnrollmentCost)
}

func assertDTOEqualCourseSection(t *testing.T, dto *dtos.CourseSectionDTO, section *entities.CourseSection) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, section.ID, dto.ID)
	require.Equal(t, section.CourseOfferingID, dto.CourseOfferingID)
	require.Equal(t, section.Name, dto.Name)
	require.Equal(t, section.Description, dto.Description)
	require.Equal(t, section.Order, dto.Order)
	require.Equal(t, string(section.Status), dto.Status)
}

func assertDTOEqualSectionModule(t *testing.T, dto *dtos.SectionModuleDTO, module *entities.SectionModule) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, module.ID, dto.ID)
	require.Equal(t, module.CourseSectionID, dto.CourseSectionID)
	require.Equal(t, module.Name, dto.Name)
	require.Equal(t, module.Description, dto.Description)
	require.Equal(t, string(module.ContentType), dto.ContentType)
	require.Equal(t, string(module.ContentStatus), dto.ContentStatus)
	require.Equal(t, module.Order, dto.Order)
	if module.ContentID != nil {
		require.Equal(t, module.ContentID, dto.ContentID)
	}
}

func assertDTOEqualCourseOfferingInstructor(t *testing.T, dto *dtos.CourseOfferingInstructorDTO, instructor *entities.CourseOfferingInstructor) {
	t.Helper()
	require.NotNil(t, dto)
	require.Equal(t, instructor.ID, dto.ID)
	require.Equal(t, instructor.CourseOfferingID, dto.CourseOfferingID)
	require.Equal(t, instructor.InstructorID, dto.InstructorID)
	require.Equal(t, instructor.InstructorUsername, dto.InstructorUsername)
}

