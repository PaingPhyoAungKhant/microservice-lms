package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
)

type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

type EnrollmentQuery struct {
	SearchQuery        *string
	StudentID          *string
	CourseID           *string
	CourseOfferingID   *string
	Status             *valueobjects.EnrollmentStatus
	Limit              *int
	Offset             *int
	SortColumn         *string
	SortDirection      *SortDirection
}

type EnrollmentQueryResult struct {
	Enrollments []*entities.Enrollment
	Total      int
}

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entities.Enrollment) error
	FindByID(ctx context.Context, id string) (*entities.Enrollment, error)
	Update(ctx context.Context, enrollment *entities.Enrollment) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, query EnrollmentQuery) (*EnrollmentQueryResult, error)
	UpdateStudentUsername(ctx context.Context, studentID, username string) error
	UpdateCourseName(ctx context.Context, courseID, courseName string) error
	UpdateCourseOfferingName(ctx context.Context, courseOfferingID, courseOfferingName string) error
}

