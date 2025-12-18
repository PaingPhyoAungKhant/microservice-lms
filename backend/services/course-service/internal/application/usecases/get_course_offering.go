package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type GetCourseOfferingUseCase struct {
	offeringRepo   repositories.CourseOfferingRepository
	courseRepo     repositories.CourseRepository
	instructorRepo repositories.CourseOfferingInstructorRepository
	sectionRepo    repositories.CourseSectionRepository
	moduleRepo     repositories.SectionModuleRepository
}

func NewGetCourseOfferingUseCase(
	offeringRepo repositories.CourseOfferingRepository,
	courseRepo repositories.CourseRepository,
	instructorRepo repositories.CourseOfferingInstructorRepository,
	sectionRepo repositories.CourseSectionRepository,
	moduleRepo repositories.SectionModuleRepository,
) *GetCourseOfferingUseCase {
	return &GetCourseOfferingUseCase{
		offeringRepo:   offeringRepo,
		courseRepo:     courseRepo,
		instructorRepo: instructorRepo,
		sectionRepo:    sectionRepo,
		moduleRepo:     moduleRepo,
	}
}

func (uc *GetCourseOfferingUseCase) Execute(ctx context.Context, offeringID string) (*dtos.CourseOfferingDetailDTO, error) {
	offering, err := uc.offeringRepo.FindByID(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, ErrCourseOfferingNotFound
	}

	var offeringDTO dtos.CourseOfferingDTO
	offeringDTO.FromEntity(offering)
	
	course, err := uc.courseRepo.FindByID(ctx, offering.CourseID)
	if err == nil && course != nil {
		courseName := course.Name
		offeringDTO.CourseName = &courseName
	}

	instructors, err := uc.instructorRepo.FindByOfferingID(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	instructorDTOs := []dtos.CourseOfferingInstructorDTO{}
	for _, instructor := range instructors {
		var instructorDTO dtos.CourseOfferingInstructorDTO
		instructorDTO.FromEntity(instructor)
		instructorDTOs = append(instructorDTOs, instructorDTO)
	}

	sections, err := uc.sectionRepo.FindByOfferingID(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	sectionDetails := []dtos.CourseSectionDetailDTO{}
	for _, section := range sections {
		var sectionDTO dtos.CourseSectionDTO
		sectionDTO.FromEntity(section)

		modules, err := uc.moduleRepo.FindBySectionID(ctx, section.ID)
		if err != nil {
			return nil, err
		}

		moduleDTOs := []dtos.SectionModuleDTO{}
		for _, module := range modules {
			var moduleDTO dtos.SectionModuleDTO
			moduleDTO.FromEntity(module)
			moduleDTOs = append(moduleDTOs, moduleDTO)
		}

		sectionDetails = append(sectionDetails, dtos.CourseSectionDetailDTO{
			CourseSectionDTO: sectionDTO,
			Modules:          moduleDTOs,
		})
	}

	return &dtos.CourseOfferingDetailDTO{
		CourseOfferingDTO: offeringDTO,
		Instructors:       instructorDTOs,
		Sections:          sectionDetails,
	}, nil
}

