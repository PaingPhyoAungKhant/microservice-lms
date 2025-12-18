package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type GetCourseWithDetailsUseCase struct {
	courseRepo           repositories.CourseRepository
	offeringRepo         repositories.CourseOfferingRepository
	instructorRepo       repositories.CourseOfferingInstructorRepository
	sectionRepo          repositories.CourseSectionRepository
	moduleRepo           repositories.SectionModuleRepository
	apiGatewayURL        string
}

func NewGetCourseWithDetailsUseCase(
	courseRepo repositories.CourseRepository,
	offeringRepo repositories.CourseOfferingRepository,
	instructorRepo repositories.CourseOfferingInstructorRepository,
	sectionRepo repositories.CourseSectionRepository,
	moduleRepo repositories.SectionModuleRepository,
	apiGatewayURL string,
) *GetCourseWithDetailsUseCase {
	return &GetCourseWithDetailsUseCase{
		courseRepo:     courseRepo,
		offeringRepo:    offeringRepo,
		instructorRepo:  instructorRepo,
		sectionRepo:    sectionRepo,
		moduleRepo:     moduleRepo,
		apiGatewayURL:  apiGatewayURL,
	}
}

func (uc *GetCourseWithDetailsUseCase) Execute(ctx context.Context, courseID string) (*dtos.CourseDetailDTO, error) {
	course, err := uc.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	var courseDTO dtos.CourseDTO
	courseDTO.FromEntity(course, uc.apiGatewayURL)

	offerings, err := uc.offeringRepo.FindByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	offeringDetails := []dtos.CourseOfferingDetailDTO{}
	for _, offering := range offerings {
		var offeringDTO dtos.CourseOfferingDTO
		offeringDTO.FromEntity(offering)

		instructors, err := uc.instructorRepo.FindByOfferingID(ctx, offering.ID)
		if err != nil {
			return nil, err
		}

		instructorDTOs := []dtos.CourseOfferingInstructorDTO{}
		for _, instructor := range instructors {
			var instructorDTO dtos.CourseOfferingInstructorDTO
			instructorDTO.FromEntity(instructor)
			instructorDTOs = append(instructorDTOs, instructorDTO)
		}

		sections, err := uc.sectionRepo.FindByOfferingID(ctx, offering.ID)
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

		offeringDetails = append(offeringDetails, dtos.CourseOfferingDetailDTO{
			CourseOfferingDTO: offeringDTO,
			Instructors:       instructorDTOs,
			Sections:          sectionDetails,
		})
	}

	return &dtos.CourseDetailDTO{
		Course:    courseDTO,
		Offerings: offeringDetails,
	}, nil
}

