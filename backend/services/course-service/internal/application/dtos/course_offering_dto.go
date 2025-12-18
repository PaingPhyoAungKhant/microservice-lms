package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseOfferingDTO struct {
	ID             string     `json:"id"`
	CourseID       string     `json:"course_id"`
	CourseName     *string    `json:"course_name,omitempty"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	OfferingType   string     `json:"offering_type"`
	Status         string     `json:"status"`
	Duration       *string    `json:"duration,omitempty"`
	ClassTime      *string    `json:"class_time,omitempty"`
	EnrollmentCost float64    `json:"enrollment_cost"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (d *CourseOfferingDTO) FromEntity(offering *entities.CourseOffering) {
	d.ID = offering.ID
	d.CourseID = offering.CourseID
	d.Name = offering.Name
	d.Description = offering.Description
	d.OfferingType = string(offering.OfferingType)
	d.Status = string(offering.Status)
	d.Duration = offering.Duration
	d.ClassTime = offering.ClassTime
	d.EnrollmentCost = offering.EnrollmentCost
	d.CreatedAt = offering.CreatedAt
	d.UpdatedAt = offering.UpdatedAt
}

type CreateCourseOfferingInput struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	OfferingType   string  `json:"offering_type" binding:"required"`
	Duration       *string `json:"duration,omitempty"`
	ClassTime      *string `json:"class_time,omitempty"`
	EnrollmentCost float64 `json:"enrollment_cost" binding:"required"`
}

type UpdateCourseOfferingInput struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	OfferingType   string  `json:"offering_type" binding:"required"`
	Duration       *string `json:"duration,omitempty"`
	ClassTime      *string `json:"class_time,omitempty"`
	EnrollmentCost float64 `json:"enrollment_cost" binding:"required"`
	Status         *string `json:"status,omitempty"`
}

type AssignInstructorInput struct {
	InstructorID       string `json:"instructor_id" binding:"required"`
	InstructorUsername string `json:"instructor_username" binding:"required"`
}

type CourseOfferingInstructorDTO struct {
	ID                 string    `json:"id"`
	CourseOfferingID   string    `json:"course_offering_id"`
	InstructorID       string    `json:"instructor_id"`
	InstructorUsername string    `json:"instructor_username"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (d *CourseOfferingInstructorDTO) FromEntity(instructor *entities.CourseOfferingInstructor) {
	d.ID = instructor.ID
	d.CourseOfferingID = instructor.CourseOfferingID
	d.InstructorID = instructor.InstructorID
	d.InstructorUsername = instructor.InstructorUsername
	d.CreatedAt = instructor.CreatedAt
	d.UpdatedAt = instructor.UpdatedAt
}

