package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
)

type EnrollmentDTO struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	StudentUsername    string    `json:"student_username"`
	CourseID           string    `json:"course_id"`
	CourseName         string    `json:"course_name"`
	CourseOfferingID   string    `json:"course_offering_id"`
	CourseOfferingName string    `json:"course_offering_name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (dto *EnrollmentDTO) FromEntity(enrollment *entities.Enrollment) {
	dto.ID = enrollment.ID
	dto.StudentID = enrollment.StudentID
	dto.StudentUsername = enrollment.StudentUsername
	dto.CourseID = enrollment.CourseID
	dto.CourseName = enrollment.CourseName
	dto.CourseOfferingID = enrollment.CourseOfferingID
	dto.CourseOfferingName = enrollment.CourseOfferingName
	dto.Status = enrollment.Status.String()
	dto.CreatedAt = enrollment.CreatedAt
	dto.UpdatedAt = enrollment.UpdatedAt
}

