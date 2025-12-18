package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
)

type Enrollment struct {
	ID                  string
	StudentID           string
	StudentUsername     string
	CourseID            string
	CourseName          string
	CourseOfferingID    string
	CourseOfferingName  string
	Status              valueobjects.EnrollmentStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewEnrollment(
	studentID string,
	studentUsername string,
	courseID string,
	courseName string,
	courseOfferingID string,
	courseOfferingName string,
) *Enrollment {
	now := time.Now().UTC()
	return &Enrollment{
		ID:                 uuid.New().String(),
		StudentID:         studentID,
		StudentUsername:    studentUsername,
		CourseID:           courseID,
		CourseName:         courseName,
		CourseOfferingID:   courseOfferingID,
		CourseOfferingName: courseOfferingName,
		Status:             valueobjects.EnrollmentStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

