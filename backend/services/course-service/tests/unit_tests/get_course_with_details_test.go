package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCourseWithDetails_Success(t *testing.T) {
	courseRepo := new(mocks.MockCourseRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	moduleRepo := new(mocks.MockSectionModuleRepository)

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     courseID,
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	instructorID := "instructor-123"
	instructor := &entities.CourseOfferingInstructor{
		ID:                 "instructor-assignment-123",
		CourseOfferingID:   offeringID,
		InstructorID:       instructorID,
		InstructorUsername: "instructor_user",
	}

	sectionID := "section-123"
	section := &entities.CourseSection{
		ID:               sectionID,
		CourseOfferingID: offeringID,
		Name:             "Introduction",
		Description:      "Introduction section",
		Order:            1,
		Status:           entities.SectionStatusDraft,
	}

	moduleID := "module-123"
	module := &entities.SectionModule{
		ID:              moduleID,
		CourseSectionID: sectionID,
		Name:            "Module 1",
		Description:     "First module",
		ContentType:     entities.ContentTypeZoom,
		ContentStatus:   entities.ContentStatusDraft,
		Order:           1,
	}

	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("FindByCourseID", mock.Anything, courseID).Return([]*entities.CourseOffering{offering}, nil).Once()
	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return([]*entities.CourseOfferingInstructor{instructor}, nil).Once()
	sectionRepo.On("FindByOfferingID", mock.Anything, offeringID).Return([]*entities.CourseSection{section}, nil).Once()
	moduleRepo.On("FindBySectionID", mock.Anything, sectionID).Return([]*entities.SectionModule{module}, nil).Once()

	uc := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := uc.Execute(context.Background(), courseID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, courseID, result.Course.ID)
	require.Len(t, result.Offerings, 1)
	require.Equal(t, offeringID, result.Offerings[0].ID)
	require.Len(t, result.Offerings[0].Instructors, 1)
	require.Equal(t, instructorID, result.Offerings[0].Instructors[0].InstructorID)
	require.Len(t, result.Offerings[0].Sections, 1)
	require.Equal(t, sectionID, result.Offerings[0].Sections[0].ID)
	require.Len(t, result.Offerings[0].Sections[0].Modules, 1)
	require.Equal(t, moduleID, result.Offerings[0].Sections[0].Modules[0].ID)

	courseRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	instructorRepo.AssertExpectations(t)
	sectionRepo.AssertExpectations(t)
	moduleRepo.AssertExpectations(t)
}

func TestGetCourseWithDetails_CourseNotFound(t *testing.T) {
	courseRepo := new(mocks.MockCourseRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	moduleRepo := new(mocks.MockSectionModuleRepository)

	courseID := "course-123"

	courseRepo.On("FindByID", mock.Anything, courseID).Return((*entities.Course)(nil), nil).Once()

	uc := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	_, err := uc.Execute(context.Background(), courseID)
	require.Error(t, err)
	require.ErrorIs(t, err, usecases.ErrCourseNotFound)

	offeringRepo.AssertNotCalled(t, "FindByCourseID", mock.Anything, mock.Anything)
	courseRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
}

func TestGetCourseWithDetails_EmptyOfferings(t *testing.T) {
	courseRepo := new(mocks.MockCourseRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	moduleRepo := new(mocks.MockSectionModuleRepository)

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("FindByCourseID", mock.Anything, courseID).Return([]*entities.CourseOffering{}, nil).Once()

	uc := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := uc.Execute(context.Background(), courseID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, courseID, result.Course.ID)
	require.Len(t, result.Offerings, 0)

	courseRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
}

func TestGetCourseWithDetails_OfferingWithNoSections(t *testing.T) {
	courseRepo := new(mocks.MockCourseRepository)
	offeringRepo := new(mocks.MockCourseOfferingRepository)
	instructorRepo := new(mocks.MockCourseOfferingInstructorRepository)
	sectionRepo := new(mocks.MockCourseSectionRepository)
	moduleRepo := new(mocks.MockSectionModuleRepository)

	courseID := "course-123"
	course := &entities.Course{
		ID:          courseID,
		Name:        "Test Course",
		Description: "Test Description",
	}

	offeringID := "offering-123"
	offering := &entities.CourseOffering{
		ID:           offeringID,
		CourseID:     courseID,
		Name:         "Spring 2024",
		OfferingType: entities.OfferingTypeOnline,
		Status:       entities.OfferingStatusPending,
	}

	courseRepo.On("FindByID", mock.Anything, courseID).Return(course, nil).Once()
	offeringRepo.On("FindByCourseID", mock.Anything, courseID).Return([]*entities.CourseOffering{offering}, nil).Once()
	instructorRepo.On("FindByOfferingID", mock.Anything, offeringID).Return([]*entities.CourseOfferingInstructor{}, nil).Once()
	sectionRepo.On("FindByOfferingID", mock.Anything, offeringID).Return([]*entities.CourseSection{}, nil).Once()

	uc := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		"http://localhost:3000",
	)

	result, err := uc.Execute(context.Background(), courseID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Offerings, 1)
	require.Len(t, result.Offerings[0].Instructors, 0)
	require.Len(t, result.Offerings[0].Sections, 0)

	courseRepo.AssertExpectations(t)
	offeringRepo.AssertExpectations(t)
	instructorRepo.AssertExpectations(t)
	sectionRepo.AssertExpectations(t)
	moduleRepo.AssertExpectations(t)
}

