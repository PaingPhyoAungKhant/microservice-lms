package dtos

type CourseDetailDTO struct {
	Course     CourseDTO                `json:"course"`
	Offerings  []CourseOfferingDetailDTO `json:"offerings"`
}

type CourseOfferingDetailDTO struct {
	CourseOfferingDTO
	Instructors []CourseOfferingInstructorDTO `json:"instructors"`
	Sections    []CourseSectionDetailDTO       `json:"sections"`
}

type CourseSectionDetailDTO struct {
	CourseSectionDTO
	Modules []SectionModuleDTO `json:"modules"`
}

